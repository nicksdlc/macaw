package context

import (
	"testing"

	"github.com/nicksdlc/macaw/config"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnErrIfNoMockMentioned(t *testing.T) {
	cfg := config.Configuration{}

	_, err := BuildContext(cfg)

	assert.Error(t, err, "Mock profile is not defined")
}

func TestShouldReturnErrIfNotSupportedProtocolProvided(t *testing.T) {
	cfg := config.Configuration{
		Mock: "NOT_SUPPORTED",
	}

	_, err := BuildContext(cfg)

	assert.Error(t, err, "Not supported protocol to mock")
}

func TestShouldReturnErrIfNocommunicatorConfigurationProvided(t *testing.T) {
	cfg := config.Configuration{
		Mock: "HTTP",
	}

	_, err := BuildContext(cfg)

	assert.Error(t, err, "communicator configuration is missing")
}
