package context

import (
	"fmt"

	"github.com/nicksdlc/macaw/communicators"
	"github.com/nicksdlc/macaw/config"
)

// communicatorBuilder is a function to build communicator from configuration file
type communicatorBuilder func(*config.Configuration) (communicators.Communicator, error)

var communicatorBuilders = make(map[string]communicatorBuilder)

func init() {
	communicatorBuilders["HTTP"] = buildHTTPCommunicator
	communicatorBuilders["RabbitMQ"] = buildRMQCommunicator

}

// BuildCommunicator is a factory to build the communicator for the context
func BuildCommunicator(cfg *config.Configuration) (communicators.Communicator, error) {
	builder, ok := communicatorBuilders[cfg.Mock]
	if !ok {
		return nil, fmt.Errorf("Not supported protocol to mock")
	}

	return builder(cfg)
}

func buildHTTPCommunicator(cfg *config.Configuration) (communicators.Communicator, error) {
	if cfg.HTTP == (config.HTTP{}) {
		return nil, fmt.Errorf("communicator configuration is missing")
	}
	return communicators.NewHTTPCommunicator(cfg.HTTP.Host, cfg.HTTP.Port, nil), nil
}

func buildRMQCommunicator(cfg *config.Configuration) (communicators.Communicator, error) {
	if cfg.Rabbit == (config.RabbitMQ{}) {
		return nil, fmt.Errorf("communicator configuration is missing")
	}
	rmqConnectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.Rabbit.User, cfg.Rabbit.Password, cfg.Rabbit.Host, cfg.Rabbit.Port)
	return communicators.NewRMQExchangeCommunicator(
		rmqConnectionString,
		cfg.Rabbit.ConnectionRetry,
		cfg.Rabbit.ResponseExchange,
		cfg.Rabbit.RequestQueue,
		cfg.Rabbit.ResponseQueue), nil
}
