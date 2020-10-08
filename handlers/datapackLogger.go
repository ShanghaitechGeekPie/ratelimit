package handlers

import "log"

func DatapackLogger(channel chan Datapack) {
	for {
		_, ok := <-channel
		if !ok {
			log.Println("DatapackLogger Disconnected")
			break
		}
	}
}
