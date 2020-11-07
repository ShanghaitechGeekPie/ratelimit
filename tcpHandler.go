package main

import (
	"log"
	"net"
	"strings"
	"sync"
)

var BucketList []*Bucket

func addrStringToIP(addrString string) net.IP {
	return net.ParseIP(addrString[:strings.LastIndexByte(addrString, ':')])
}

func policyResult(config *ServerConfig, ip net.IP) int {
	for index, v := range config.Policy {
		if v.Source.Contains(ip) {
			return index
		}
	}
	panic(ip.String() + "is not included in the rule")
}

func handleCopy(connFrom net.Conn, connTo net.Conn, bucket *Bucket, wg *sync.WaitGroup) {
	defer wg.Done()
	buffer := make([]byte, bufferSize)
	for {
		readLength, err := connFrom.Read(buffer)
		bucket.Wait(int64(readLength))
		if err != nil || readLength == 0 {
			break
		}
		if readLength == bufferSize {
			_, err = connTo.Write(buffer)
		} else {
			_, err = connTo.Write(buffer[:readLength])
		}
		if err != nil {
			break
		}
	}
}

func handleForward(config *ServerConfig, connectionToLocal net.Conn) {
	policyID := policyResult(config, addrStringToIP(connectionToLocal.RemoteAddr().String()))
	connectionToDestination, err := net.Dial("tcp", config.Policy[policyID].Destination)
	if err != nil {
		_ = connectionToLocal.Close()
	} else {
		var wg sync.WaitGroup
		wg.Add(2)
		go handleCopy(connectionToLocal, connectionToDestination, BucketList[policyID], &wg)
		go handleCopy(connectionToDestination, connectionToLocal, BucketList[policyID], &wg)
		wg.Wait()
		_ = connectionToLocal.Close()
		_ = connectionToDestination.Close()
	}
}

func HandleTCP(config *ServerConfig) {
	lis, err := net.Listen("tcp", config.Listen)
	if err != nil {
		log.Fatal(err)
	}
	//defer lis.Close()
	for {
		connectionToLocal, err := lis.Accept()
		if err == nil {
			log.Print(connectionToLocal.)
			go handleForward(config, connectionToLocal)
		}
	}
}
