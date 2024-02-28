package admin

import (
	"github.com/nicksdlc/macaw/model"
	"github.com/nicksdlc/macaw/prototype"
)

type MockCommunicator struct {
	Prototypes []prototype.MessagePrototype
}

// RespondWith is a mock method
func (m *MockCommunicator) RespondWith(response []prototype.MessagePrototype) {
	m.Prototypes = response
}

// RequestWith is a mock method
func (m *MockCommunicator) RequestWith(request []prototype.MessagePrototype) {
}

// GetResponses is a mock method
func (m *MockCommunicator) GetResponses() []prototype.MessagePrototype {
	return m.Prototypes
}

// UpdateResponse is a mock method
func (m *MockCommunicator) UpdateResponse(response prototype.MessagePrototype) {
	for i, resp := range m.Prototypes {
		if resp.Alias == response.Alias {
			m.Prototypes[i] = response
		}
	}
}

// Close is a mock method
func (m *MockCommunicator) Close() error {
	return nil
}

// PostAndListen is a mock method
func (m *MockCommunicator) PostAndListen() (chan model.ResponseMessage, error) {
	return nil, nil
}

// ConsumeMediateReplyWithResponse is a mock method
func (m *MockCommunicator) ConsumeMediateReplyWithResponse() {
}
