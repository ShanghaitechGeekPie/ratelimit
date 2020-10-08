package handlers

import "log"

func BucketEmptyLogger(channel chan int) {
	for {
		ID, err := <-channel
		if !err {
			log.Println(ID, " exhausted")
		} else {
			log.Println(ID, "BucketEmptyLogger Exited")
			break
		}
	}
}
