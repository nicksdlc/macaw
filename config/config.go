package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Configuration of the macaw
type Configuration struct {
	Admin          Admin
	ConnectThrough string
	DumpMetrics    DumpMetrics
	Rabbit         RabbitMQ
	HTTP           HTTP
	Mode           string
	Responses      []Response
	Requests       []Request
}

// Admin configuration
type Admin struct {
	Enabled bool
	Port    uint16
}

// Request configuration
type Request struct {
	Alias   string
	To      string
	From    string
	Type    string
	Body    Body
	Options Options
}

// DumpMetrics configuration
type DumpMetrics struct {
	Enabled   bool
	Frequency int
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

// ReadConfigFromArg reads the configuration file path from the command line
func ReadConfigFromArg() Configuration {
	pflag.StringP("config", "c", "config.yml", "Path to configuration file")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	cfg := viper.GetString("config")
	cfg = strings.TrimSuffix(cfg, filepath.Ext(cfg))

	fmt.Printf("Starting with config file: %s. Use --config, -c to specify another file\n", cfg)
	return Read(cfg)
}
