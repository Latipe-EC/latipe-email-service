package main

import (
	"email-service/config"
	"email-service/service"
	message "email-service/worker"
	"log"
)

func main() {
	log.Printf("--= Init Email Worker =--")
	globalCfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Init worker was failed cause: %v", err)
	}

	gmailSender := service.NewGmailSenderEmail(globalCfg)

	gmailWorker := message.NewConsumerEmailWorker(globalCfg, gmailSender)
	log.Printf("Create connection to queue")

	listenQueue := make(chan bool)

	go func() {
		gmailWorker.ListenMessageQueue()
	}()
	<-listenQueue

}
