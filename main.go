package main

import (
	"./handlers"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type configuration struct {
	Data []handlers.PortConfig
}

func main() {
	var wg sync.WaitGroup

	dataPackChannel := make(chan handlers.Datapack)
	defer close(dataPackChannel)
	go handlers.DatapackLogger(dataPackChannel)

	bucketEmptyChannel := make(chan int)
	defer close(bucketEmptyChannel)
	go handlers.BucketEmptyLogger(bucketEmptyChannel)

	var closeMu sync.RWMutex
	var closeBool = false
	closedSignal := handlers.Signal{Mu: &closeMu, Value: &closeBool}

	file, _ := os.Open("config.json")

	decoder := json.NewDecoder(file)
	conf := configuration{}
	cfgerr := decoder.Decode(&conf)
	if cfgerr != nil {
		fmt.Println("Error:", cfgerr)
	}
	fmt.Println(conf.Data)
	var muSlice = make([]sync.RWMutex, len(conf.Data))
	var valSlice = make([]handlers.Value1, len(conf.Data))
	var sig1Slice = make([]handlers.Signal1, len(conf.Data))

	for index := range sig1Slice {
		valSlice[index].CanModify = true
		sig1Slice[index].Mu = &muSlice[index]
		sig1Slice[index].Value = &valSlice[index]
	}
	for index, value := range conf.Data {
		wg.Add(1)
		go handlers.PortHandler(value, dataPackChannel, closedSignal, bucketEmptyChannel, sig1Slice[index], &wg)
	}
	go func(){
		time.Sleep(10*time.Second)
		*closedSignal.Value=true
		log.Println("Go Closed")
	}()
	wg.Wait()
}
