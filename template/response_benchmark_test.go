package template

import (
	"testing"
)

func BenchmarkShouldReplaceWithRandomNumber(b *testing.B) {
	resp := NewResponse(
		`{
			number: {{.Number}}, 
			string: {{.String}}, 
			date: {{.Date}}, 
			incremental: {{.Number "incremental"}},
			list: {{.List "[{{.Number}}, {{.Number}}]" 3}}
		}`, 2, nil)

	for i := 0; i < b.N; i++ {
		resp.Create()
	}
}

func BenchmarkWithNumber(b *testing.B) {
	resp := NewResponse(
		`{
			number: {{.Number "between" "10" "10"}}
		}`, 1, nil)

	for i := 0; i < b.N; i++ {
		resp.Create()
	}
}

func BenchmarkWithNumberAndString(b *testing.B) {
	resp := NewResponse(
		`{
			number: {{.Number "between" "10" "10"}}
			string: {{.String}}, 
		}`, 1, nil)

	for i := 0; i < b.N; i++ {
		resp.Create()
	}
}

func BenchmarkWithNumberStringAndList(b *testing.B) {
	resp := NewResponse(
		`{
			number: {{.Number "between" "10" "10"}}
			string: {{.String}}, 
			list: {{.List "[{{.Number}}, {{.Number}}]" 3}}
		}`, 1, nil)

	for i := 0; i < b.N; i++ {
		resp.Create()
	}
}
