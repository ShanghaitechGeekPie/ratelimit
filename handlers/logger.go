package handlers

import "log"

func logger(channel chan string) {
	logToSend := ""
	for {
		newLog, ok := <-channel
		if ok {
			log.Println(newLog)
			logToSend += newLog + "\n"
		} else {
			log.Println("LogClosed")
			break
		}
	}
}
