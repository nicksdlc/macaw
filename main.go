package main

import (
	"fmt"
	"log"
	"macaw/config"
	"macaw/connectors"
	"macaw/receiver"
	"macaw/sender"
)

func main() {
	cfg := readConfig()

	rmqConnectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.Rabbit.User, cfg.Rabbit.Password, cfg.Rabbit.Host, cfg.Rabbit.Port)
	rc := connectors.NewRMQExchangeConnector(rmqConnectionString, cfg.Rabbit.ResponseExchange, cfg.Rabbit.RequestQueue, cfg.Rabbit.ResponseQueue)
	defer rc.Close()

	if cfg.Mode == "receiver" {
		listener := receiver.NewRMQReceiver(rc, cfg.Response)
		listener.Listen()

		var forever chan struct{}

		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
		<-forever
	}

	if cfg.Mode == "sender" {
		sender := sender.NewRMQSender(rc, cfg.Request)
		sender.Send()
	}

}

func readConfig() config.Configuration {
	return config.Read()
}
