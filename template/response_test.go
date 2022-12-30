package template

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldReplaceWithRandomNumber(t *testing.T) {
	resp := NewResponse(
		`{
			number: {{.Number}}, 
			string: {{.String}}, 
			date: {{.Date}}, 
			incremental: {{.Number "incremental"}},
			list: {{.List "[{{.Number}}, {{.Number}}]" 3}}
		}`, 2, nil)

	first := resp.Create()
	second := resp.Create()

	assert.NotEqual(t, first, second)
	assert.Contains(t, first, "incremental: 1")
	assert.Contains(t, second, "incremental: 2")
	_, err := json.Marshal([]byte(first))
	assert.NoError(t, err)
}

func TestShouldReplaceStringTypes(t *testing.T) {
	resp := NewResponse(`{ 
			oneOf: {{.String "variant" "one" "another"}}
		}`, 1, nil)

	result := resp.Create()

	if strings.Contains(result, "oneOf: one") || strings.Contains(result, "oneOf: another") {
		fmt.Println(result)
		t.Log("Passed")
	} else {
		t.Logf("Want: oneOf: one or oneOf: another, got %s", result)
		t.Fail()
	}
}

func TestShouldReplaceWithCorrectAmount(t *testing.T) {
	resp := NewResponse(
		`{
			number: {{.Amount}}
		}`, 2, nil)

	result := resp.Create()

	assert.Contains(t, result, "number: 2")
}

func TestShouldReplaceFromRequset(t *testing.T) {
	req := Serialize([]byte(`{ "id": 123, "task": "DoSomething" }`))
	resp := NewResponse(
		`{
			incremental: {{.Number "incremental"}},
			request: {{.FromRequest "id"}}
		}`, 1, &req)

	result := resp.Create()

	assert.Contains(t, result, "incremental: 1")
	assert.Contains(t, result, "request: 123")
}

func TestShouldReplaceFromRequsetWithHierarchy(t *testing.T) {
	req := Serialize([]byte(`{ "payload": {"id": 123, "field":"test"}, "task": "DoSomething" }`))
	resp := NewResponse(
		`{
			incremental: {{.Number "incremental"}},
			request: {{.FromRequest "payload/field"}}
		}`, 1, &req)

	result := resp.Create()

	assert.Contains(t, result, "incremental: 1")
	assert.Contains(t, result, "request: test")
}
