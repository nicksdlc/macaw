package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Configuration of the macaw
type Configuration struct {
	ConnectThrough string
	Rabbit         RabbitMQ
	HTTP           HTTP
	Mode           string
	Responses      []Response
	Request        Request
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

// Read configuration from file in the same directory as executable
func Read(name string) Configuration {
	viper.SetConfigName(name)
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
