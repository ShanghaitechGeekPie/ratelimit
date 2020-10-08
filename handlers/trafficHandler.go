package handlers

import (
	"../ratelimit"
	"io"
	"net"
	"sync"
	"time"
)

type PortConfig struct {
	Bandwidth  int64
	LocalAddr  string
	RemoteAddr string
	ID         int
}

type Datapack struct {
	size int
	id   int
}

//the time gap that goroutine send a signal showing that the bucket is empty
const BucketRefreshWaitTime = int64(500 * time.Millisecond)

func checkSignal1Available(currentTick int64) func(value *Value1) bool {
	return func(value *Value1) bool {
		return value.CanModify && value.LatestModified-currentTick > BucketRefreshWaitTime
	}
}

func writeSignal1Available(currentTick int64) func(value *Value1) {
	return func(value *Value1) {
		value.LatestModified = currentTick
	}
}

func bucketEmptyTrigger(bucketEmptyInfoCanWrite Signal1, id int, channel chan int) func(LastTick int64) {
	return func(LastTick int64) {
		if bucketEmptyInfoCanWrite.ReadSignal1(checkSignal1Available(LastTick)) {
			bucketEmptyInfoCanWrite.WriteSignal1(writeSignal1Available(LastTick))
			channel <- id
		}
	}
}

func trafficHandler(local net.Conn, remote net.Conn, wg *sync.WaitGroup, bucket *ratelimit.Bucket, channel chan string, byteChannel chan Datapack, empty chan int, signalBucketEmpty Signal1, id int) {
	defer wg.Done()
	reader := Reader(local, bucket)
	writer := Writer(remote, bucket)
	_, err := copyBuffer(writer, reader, bucketEmptyTrigger(signalBucketEmpty, id, empty), byteChannel, id)
	if err != nil && err != io.EOF {
		channel <- err.Error()
	}
	channel <- "exited."
	remote.Close()
}

func twoWayConnHandler(localConn net.Conn, cfg PortConfig, bucket *ratelimit.Bucket, byteChannel chan Datapack, logChannel chan string, empty chan int, signalBucketEmpty Signal1, wg *sync.WaitGroup) {
	remoteConn, err := net.Dial("tcp", cfg.RemoteAddr)
	if err != nil {
		logChannel <- err.Error()
		wg.Done()
		wg.Done()
	} else {
		go trafficHandler(localConn, remoteConn, wg, bucket, logChannel, byteChannel, empty, signalBucketEmpty, cfg.ID)
		go trafficHandler(remoteConn, localConn, wg, bucket, logChannel, byteChannel, empty, signalBucketEmpty, cfg.ID)
	}
}


// true => closed
// false => running
func validateCloseSignal(signal *bool) bool {
	return *signal
}

//bucketEmptyChannel provides the ID of Bucket that is empty
//byteChannel provides the amount of data that was consumed by the port
//closed provides whether the main process want to end up the goroutine

func PortHandler(cfg PortConfig, byteChannel chan Datapack, closed Signal, bucketEmptyChannel chan int, signalBucketEmpty Signal1, wgMain *sync.WaitGroup) {
	defer wgMain.Done()

	var logChannel = make(chan string)
	defer close(logChannel)
	go logger(logChannel)

	lis, err := net.Listen("tcp", cfg.LocalAddr)
	if err != nil {
		logChannel <- err.Error()
	} else {
		defer lis.Close()
		var wg sync.WaitGroup
		bucket := ratelimit.NewBucketWithQuantum(time.Second, 100*ratelimit.Megabyte.GetByte(), ratelimit.Megabyte.GetByte()*2)
		for {
			localConn, err := lis.Accept()
			if err != nil {
				logChannel <- err.Error()
			} else if closed.ReadSignal(validateCloseSignal){
				localConn.Close()
				break
			} else {
				wg.Add(2)
				twoWayConnHandler(localConn, cfg, bucket, byteChannel, logChannel, bucketEmptyChannel, signalBucketEmpty, &wg)
			}
		}
		wg.Wait()
	}
}
