package main

import (
	"email-service/config"
	"email-service/service"
	message "email-service/worker"
	"log"
	"sync"
)

func main() {
	log.Printf("--= Init Email Worker =--")
	globalCfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Init worker was failed cause: %v", err)
	}

	gmailSender := service.NewGmailSenderEmail(globalCfg)

	orderMessageWorker := message.NewConsumerOrderWorker(globalCfg, gmailSender)
	log.Printf("OrderMessage subscriber was created")

	userRegisterWorker := message.NewConsumerUserRegisterWorker(globalCfg, gmailSender)
	log.Printf("UserRegister subscriber was created")

	deliveryAccountWorker := message.NewConsumerDeliveryWorker(globalCfg, gmailSender)
	log.Printf("DeliveryAccount subscriber was created")

	forgotPasswordWorker := message.NewConsumerForgotPasswordWorker(globalCfg, gmailSender)
	log.Printf("forgotPassword subscriber was created")

	var wg sync.WaitGroup
	wg.Add(1)
	go userRegisterWorker.ListenMessageQueue(&wg)

	wg.Add(1)
	go orderMessageWorker.ListenMessageQueue(&wg)

	wg.Add(1)
	go deliveryAccountWorker.ListenMessageQueue(&wg)

	wg.Add(1)
	go forgotPasswordWorker.ListenMessageQueue(&wg)

	wg.Wait()

}
