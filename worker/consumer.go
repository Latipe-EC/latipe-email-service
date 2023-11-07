package message

import (
	"email-service/config"
	"email-service/data/dto"
	"email-service/service"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerOrderMessage struct {
	config      *config.Config
	emailSender service.SenderEmailService
}

func NewConsumerEmailWorker(config *config.Config, senderService service.SenderEmailService) *ConsumerOrderMessage {
	return &ConsumerOrderMessage{
		config:      config,
		emailSender: senderService,
	}
}

func (mq ConsumerOrderMessage) ListenMessageQueue() {
	conn, err := amqp.Dial(mq.config.RabbitMQ.Connection)
	failOnError(err, "Failed to connect to RabbitMQ")
	log.Printf("[%s] Comsumer has been connected", "INFO")

	channel, err := conn.Channel()
	defer channel.Close()
	defer conn.Close()

	// Khai báo một Exchange loại "direct"
	err = channel.ExchangeDeclare(
		mq.config.RabbitMQ.Exchange, // Tên Exchange
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("cannot declare exchange: %v", err)
	}

	// Tạo hàng đợi
	_, err = channel.QueueDeclare(
		mq.config.RabbitMQ.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("cannot declare queue: %v", err)
	}

	err = channel.QueueBind(
		mq.config.RabbitMQ.Queue,
		mq.config.RabbitMQ.RoutingKey,
		mq.config.RabbitMQ.Exchange,
		false,
		nil)
	if err != nil {
		log.Fatalf("cannot bind exchange: %v", err)
	}

	// declaring consumer with its properties over channel opened
	msgs, err := channel.Consume(
		mq.config.RabbitMQ.Queue,        // queue
		mq.config.RabbitMQ.ConsumerName, // consumer
		true,                            // auto ack
		false,                           // exclusive
		false,                           // no local
		false,                           // no wait
		nil,                             //args
	)
	if err != nil {
		panic(err)
	}

	// handle consumed messages from queue
	for msg := range msgs {
		log.Printf("[%s] received message from: %s", "INFO", msg.RoutingKey)

		if err := mq.orderHandler(msg); err != nil {
			log.Printf("[%s] The order creation failed cause %s", "ERROR", err)
		}

	}

	log.Printf("[%s] message queue has started", "INFO")
	log.Printf("[%s] waiting for messages...", "INFO")

}

func (mq ConsumerOrderMessage) orderHandler(msg amqp.Delivery) error {
	message := dto.EmailRequest{}

	if err := json.Unmarshal(msg.Body, &message); err != nil {
		log.Printf("[%s] Parse message to order failed cause: %s", "ERROR", err)
		return err
	}

	var err error
	switch message.Type {
	case dto.ORDER:
		err = mq.emailSender.SendOrderEmail(&message)
	case dto.REGISTER:
		err = mq.emailSender.SendForgotPassword(&message)
	case dto.FORGOT:
		err = mq.emailSender.SendRegisterEmail(&message)
	}
	if err != nil {
		log.Printf("[%s] send email was failed cause: %s", "ERROR", err)
	}

	return nil
}
