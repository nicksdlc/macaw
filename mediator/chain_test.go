package mediator

import (
	"testing"

	"github.com/nicksdlc/macaw/model"
	"github.com/stretchr/testify/assert"
)

type MockMediator struct {
	Name string
}

func (mm *MockMediator) Mediate(message model.RequestMessage, _ <-chan model.ResponseMessage) <-chan model.ResponseMessage {
	return nil
}

func TestMediatorChain_AppendShouldAppendInitialMediator(t *testing.T) {
	// Given
	firstMediator := &MockMediator{
		Name: "first",
	}
	mediatorChain := MediatorChain{}

	// When
	mediatorChain.Append(firstMediator)

	// Then
	if mediatorChain.linkedMediator == nil {
		t.Fatalf("expected a linked mediator, instead got: %v", mediatorChain.linkedMediator)
	}
}

func TestMediatorChain_AppendShouldAppendMediator(t *testing.T) {
	// Given
	firstMediator := &MockMediator{
		Name: "first",
	}
	secondMediator := &MockMediator{
		Name: "second",
	}
	mediatorChain := MediatorChain{}
	mediatorChain.Append(firstMediator)

	// When
	mediatorChain.Append(secondMediator)

	// Then
	assert.Equal(t, "second", mediatorChain.linkedMediator.next.mediator.(*MockMediator).Name)
}

func TestMediatorChain_AppendShouldAppendMoreThenTwoMediators(t *testing.T) {
	// Given
	firstMediator := &MockMediator{
		Name: "first",
	}
	secondMediator := &MockMediator{
		Name: "second",
	}
	thirdMediator := &MockMediator{
		Name: "third",
	}
	mediatorChain := MediatorChain{}
	mediatorChain.Append(firstMediator)
	mediatorChain.Append(secondMediator)

	// When
	mediatorChain.Append(thirdMediator)

	// Then
	assert.Equal(t, "third", mediatorChain.linkedMediator.next.next.mediator.(*MockMediator).Name)
}

func TestMediatorChain_PrependShouldPrependAllTheMediators(t *testing.T) {
	// Given
	firstMediator := &MockMediator{
		Name: "first",
	}
	secondMediator := &MockMediator{
		Name: "second",
	}
	thirdMediator := &MockMediator{
		Name: "third",
	}
	mediatorChain := MediatorChain{}

	// When
	mediatorChain.Prepend(thirdMediator)
	mediatorChain.Prepend(secondMediator)
	mediatorChain.Prepend(firstMediator)

	// Then
	assert.Equal(t, "first", mediatorChain.linkedMediator.mediator.(*MockMediator).Name)
	assert.Equal(t, "second", mediatorChain.linkedMediator.next.mediator.(*MockMediator).Name)
	assert.Equal(t, "third", mediatorChain.linkedMediator.next.next.mediator.(*MockMediator).Name)

}
