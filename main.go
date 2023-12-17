package main

import (
	"email-service/config"
	"email-service/service"
	message "email-service/worker"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
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
	log.Printf("ForgotPassword subscriber was created")

	paymentWorker := message.NewConsumerPaymentWorker(globalCfg, gmailSender)
	log.Printf("Payment subscriber was created")

	var wg sync.WaitGroup
	wg.Add(1)
	go userRegisterWorker.ListenMessageQueue(&wg)

	wg.Add(1)
	go orderMessageWorker.ListenMessageQueue(&wg)

	wg.Add(1)
	go deliveryAccountWorker.ListenMessageQueue(&wg)

	wg.Add(1)
	go forgotPasswordWorker.ListenMessageQueue(&wg)

	wg.Add(1)
	go paymentWorker.ListenMessageQueue(&wg)

	wg.Add(1)
	go func() {
		defer wg.Done()
		http.HandleFunc("/", wellcomeHandler)
		http.Handle("/metrics", promhttp.Handler())

		log.Printf("Init email service http://localhost:5015")
		err := http.ListenAndServe(":5015", nil)
		if err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}

	}()

	wg.Wait()
}

func wellcomeHandler(w http.ResponseWriter, r *http.Request) {
	s := struct {
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Message: "Email service was developed by TienDat",
		Version: "v0.0.1",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(s)
	if err != nil {
		return
	}
}
