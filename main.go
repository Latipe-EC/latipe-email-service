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

	workers := []struct {
		name   string
		worker message.Worker
	}{
		{"OrderMessage", message.NewConsumerOrderWorker(globalCfg, gmailSender, rabbitConn)},
		{"UserRegister", message.NewConsumerUserRegisterWorker(globalCfg, gmailSender, rabbitConn)},
		{"DeliveryAccount", message.NewConsumerDeliveryWorker(globalCfg, gmailSender, rabbitConn)},
		{"ForgotPassword", message.NewConsumerForgotPasswordWorker(globalCfg, gmailSender, rabbitConn)},
		{"Payment", message.NewConsumerPaymentWorker(globalCfg, gmailSender, rabbitConn)},
	}

	var wg sync.WaitGroup

	for _, worker := range workers {
		wg.Add(1)
		go runWithRecovery(func() { worker.worker.ListenMessageQueue(&wg) })
		log.Printf("%s subscriber was created", worker.name)
	}

	wg.Add(1)
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
