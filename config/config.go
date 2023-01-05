package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Configuration of the macaw
type Configuration struct {
	Mock     string
	Rabbit   RabbitMQ
	Mode     string
	Response Response
	Request  Request
}

// RabbitMQ configuration of RabbitMQ
type RabbitMQ struct {
	Host             string
	Port             string
	User             string
	Password         string
	ResponseExchange string
	RequestQueue     string
	ResponseQueue    string
	ConnectionRetry  Retry
}

// Response configuration
type Response struct {
	File   string
	Amount int
}

// Request configuration
type Request struct {
	File   string
	Amount int
	Delay  int
}

// Retry represents how many retries to do and with which interval
type Retry struct {
	ElapsedTime int
	Interval    int
}

// Read configuration from file
func Read() Configuration {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	var configuration Configuration

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	return configuration
}
