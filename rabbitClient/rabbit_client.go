package rabbitclient

import (
	"email-service/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func NewRabbitClientConnection(globalCfg *config.Config) *amqp.Connection {
	cfg := amqp.Config{
		Properties: amqp.Table{
			"connection_name": globalCfg.RabbitMQ.ServiceName,
		},
	}

	conn, err := amqp.DialConfig(globalCfg.RabbitMQ.Connection, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ cause:%v", err)
	}

	log.Println("Comsumer has been connected")
	return conn
}
