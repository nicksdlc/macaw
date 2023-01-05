package main

import (
	"fmt"
	"log"
	"macaw/config"
	"macaw/connectors"
	"macaw/receiver"

	"github.com/spf13/viper"
)

func main() {
	cfg := readConfig()

	rmqConnectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.Rabbit.User, cfg.Rabbit.Password, cfg.Rabbit.Host, cfg.Rabbit.Port)
	rc := connectors.NewRMQExchangeConnector(rmqConnectionString, cfg.Rabbit.ConnectionRetry, cfg.Rabbit.ResponseExchange, cfg.Rabbit.RequestQueue, cfg.Rabbit.ResponseQueue)
	defer rc.Close()

	listener := receiver.NewRMQReceiver(rc, cfg.ResponseTemplate)
	listener.Listen()

	var forever chan struct{}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func readConfig() config.Configuration {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	var configuration config.Configuration

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	return configuration
}
