package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	GmailHostConfig GmailHostConfig
	EmailTemplate   EmailTemplate
	RabbitMQ        RabbitMQ
}

type GmailHostConfig struct {
	EmailSender string
	Password    string
	StmpHost    string
	StmpPort    string
}

type EmailTemplate struct {
	OrderTemplate       string
	RegisterTemplate    string
	ForgotPassTemplate  string
	ConfirmLinkTemplate string
}

type RabbitMQ struct {
	Connection   string
	Exchange     string
	RoutingKey   string
	Queue        string
	ConsumerName string
	ProducerName string
}

// Get config path for local or docker
func getDefaultConfig() string {
	return "./config/config"
}

func NewConfig() (*Config, error) {
	config := Config{}
	path := os.Getenv("cfgPath")
	if path == "" {
		path = getDefaultConfig()
	}

	v := viper.New()

	v.SetConfigName(path)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	err := v.Unmarshal(&config)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &config, nil
}
