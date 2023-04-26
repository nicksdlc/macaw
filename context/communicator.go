package context

import (
	"fmt"

	"github.com/nicksdlc/macaw/communicators"
	"github.com/nicksdlc/macaw/config"
)

// communicatorBuilder is a function type to build communicator from configuration file
type communicatorBuilder func(*config.Configuration) (communicators.Communicator, error)

var communicatorBuilders = make(map[string]communicatorBuilder)

// init allows to register all the builders for the communicators
// once the package is imported
func init() {
	communicatorBuilders["HTTP"] = buildHTTPCommunicator
	communicatorBuilders["RabbitMQ"] = buildRMQCommunicator

}

// BuildCommunicator is a factory to build the communicator for the context
func BuildCommunicator(cfg *config.Configuration) (communicators.Communicator, error) {
	builder, ok := communicatorBuilders[cfg.ConnectThrough]
	if !ok {
		return nil, fmt.Errorf("not supported protocol to mock")
	}

	return builder(cfg)
}

func buildHTTPCommunicator(cfg *config.Configuration) (communicators.Communicator, error) {
	if cfg.HTTP == (config.HTTP{}) {
		return nil, fmt.Errorf("communicator configuration is missing")
	}
	return communicators.NewHTTPCommunicator(cfg.HTTP.Serve.Host, cfg.HTTP.Serve.Port, nil), nil
}

func buildRMQCommunicator(cfg *config.Configuration) (communicators.Communicator, error) {
	if !rabbitConfigIsValid(cfg) {
		return nil, fmt.Errorf("communicator configuration is missing")
	}
	rmqConnectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.Rabbit.User, cfg.Rabbit.Password, cfg.Rabbit.Host, cfg.Rabbit.Port)
	return communicators.NewRMQExchangeCommunicator(
		rmqConnectionString,
		cfg.Rabbit.ConnectionRetry,
		cfg.Rabbit.Exchanges,
		cfg.Rabbit.Queues), nil
}

func rabbitConfigIsValid(cfg *config.Configuration) bool {
	validHost := cfg.Rabbit.Host != "" && cfg.Rabbit.Port != ""
	validQueues := cfg.Rabbit.Queues != nil && len(cfg.Rabbit.Queues) > 0
	validExchanges := cfg.Rabbit.Exchanges != nil && len(cfg.Rabbit.Exchanges) > 0

	return validHost && validQueues && validExchanges
}
