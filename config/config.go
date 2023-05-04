package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Configuration of the macaw
type Configuration struct {
	Control        Control
	ConnectThrough string
	Rabbit         RabbitMQ
	HTTP           HTTP
	Mode           string
	Responses      []Response
	Request        Request
}

// Control configuration
type Control struct {
	Enabled bool
	OnPort  uint16
}

// Request configuration
type Request struct {
	To      string
	Body    Body
	Options Options
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
