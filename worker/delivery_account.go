package message

import (
	"email-service/config"
	"email-service/data/dto"
	"email-service/service"
	"encoding/json"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerDeliveryMessage struct {
	config      *config.Config
	emailSender service.SenderEmailService
}

func NewConsumerDeliveryWorker(config *config.Config, senderService service.SenderEmailService) *ConsumerDeliveryMessage {
	return &ConsumerDeliveryMessage{
		config:      config,
		emailSender: senderService,
	}
}

func (mq ConsumerDeliveryMessage) ListenMessageQueue(wg *sync.WaitGroup) {
	conn, err := amqp.Dial(mq.config.RabbitMQ.Connection)
	failOnError(err, "Failed to connect to RabbitMQ")
	log.Printf("[%s] [%s] Comsumer has been connected", "INFO", mq.config.RabbitMQ.DeliveryRegisterTopic.RoutingKey)

	channel, err := conn.Channel()
	defer channel.Close()
	defer conn.Close()

	// Khai báo một Exchange loại "direct"
	err = channel.ExchangeDeclare(
		mq.config.RabbitMQ.Exchange, // Tên Exchange
		"topic",
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
	q, err := channel.QueueDeclare(
		"",
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
		q.Name,
		mq.config.RabbitMQ.DeliveryRegisterTopic.RoutingKey,
		mq.config.RabbitMQ.Exchange,
		false,
		nil)
	if err != nil {
		log.Fatalf("cannot bind exchange: %v", err)
	}

	// declaring consumer with its properties over channel opened
	msgs, err := channel.Consume(
		q.Name,                          // queue
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

	log.Printf("[%s] [%s] message queue has started", "INFO", mq.config.RabbitMQ.DeliveryRegisterTopic.RoutingKey)
	log.Printf("[%s] [%s] waiting for messages...", "INFO", mq.config.RabbitMQ.DeliveryRegisterTopic.RoutingKey)

	// handle consumed messages from queue
	defer wg.Done()
	for msg := range msgs {
		log.Printf("[%s] received order message from: %s", "INFO", msg.RoutingKey)

		if err := mq.handleMessage(msg); err != nil {
			log.Printf("[%s] [%s] Handling message was failed cause %s", "ERROR", mq.config.RabbitMQ.DeliveryRegisterTopic.RoutingKey, err)
		}
	}
}

func (mq ConsumerDeliveryMessage) handleMessage(msg amqp.Delivery) error {
	message := dto.DeliveryAccountMessage{}

	if err := json.Unmarshal(msg.Body, &message); err != nil {
		log.Printf("[%s] Parse message to order failed cause: %s", "ERROR", err)
		return err
	}

	if err := mq.emailSender.SendDeliveryAccount(&message); err != nil {
		log.Printf("[%s] send email was failed cause: %s", "ERROR", err)
	}

	return nil
}
