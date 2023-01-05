package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Configuration struct {
	Mock             string
	Rabbit           RabbitMQ
	ResponseTemplate string
}

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

// Retry represents how many retries to do and with which interval
type Retry struct {
	ElapsedTime int
	Interval    int
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
