package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"text/template"
)

// Take a configuration file, create map[string]interface{}
// Build it with valid types and generate specific values
// builder only builds, generation is performed in the other class

type Number struct {
	Value int
}

func serialize(input string) map[string]interface{} {
	var output map[string]interface{}

	json.Unmarshal([]byte(input), &output)

	return output
}

func deserialize(input map[string]interface{}) string {
	output, _ := json.Marshal(input)

	return string(output)
}

func build(input string) string {
	num := Number{
		Value: 1,
	}

	resultTemplate, err := template.New("").
		Funcs(template.FuncMap{
			"sum": func(i, j int) string {
				return fmt.Sprintf("%d", i+j)
			},
		}).
		Parse(input)
	if err != nil {
		log.Printf("Failed to create template %s", err.Error())
	}

	var result bytes.Buffer
	err = resultTemplate.Execute(&result, num)

	return result.String()
}
