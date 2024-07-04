package main

import (
	"email-service/api"
	"log"
	"net/http"
	"sync"

	"email-service/config"
	rabbitclient "email-service/rabbitClient"
	"email-service/service"
	message "email-service/worker"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.Printf("--= Init Email Worker =--")

	globalCfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Init worker failed: %v", err)
	}

	rabbitConn := rabbitclient.NewRabbitClientConnection(globalCfg)
	gmailSender := service.NewGmailSenderEmail(globalCfg)

	orderCons := message.NewConsumerOrderWorker(globalCfg, gmailSender, rabbitConn)
	userRegCons := message.NewConsumerUserRegisterWorker(globalCfg, gmailSender, rabbitConn)
	deliRegCons := message.NewConsumerDeliveryWorker(globalCfg, gmailSender, rabbitConn)
	forgotCons := message.NewConsumerForgotPasswordWorker(globalCfg, gmailSender, rabbitConn)
	paymentCons := message.NewConsumerPaymentWorker(globalCfg, gmailSender, rabbitConn)

	var wg sync.WaitGroup
	wg.Add(1)
	go runWithRecovery(func() {
		orderCons.ListenMessageQueue(&wg)
	})

	wg.Add(1)
	go runWithRecovery(func() {
		userRegCons.ListenMessageQueue(&wg)

	})

	wg.Add(1)
	go runWithRecovery(func() {
		deliRegCons.ListenMessageQueue(&wg)

	})

	wg.Add(1)
	go runWithRecovery(func() {
		forgotCons.ListenMessageQueue(&wg)

	})

	wg.Add(1)
	go runWithRecovery(func() {
		paymentCons.ListenMessageQueue(&wg)
	})

	go runWithRecovery(func() { startHTTPServer(&wg) })

	wg.Wait()
}

func startHTTPServer(wg *sync.WaitGroup) {
	defer wg.Done()
	handler := api.NewHandlerApi()

	http.HandleFunc("/", handler.WelcomeHandler)
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Init email service http://localhost:5015")
	err := http.ListenAndServe(":5015", nil)
	if err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

func runWithRecovery(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()
	fn()
}
