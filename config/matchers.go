package config

type Matcher struct {
	// Type - matcher type
	Type string

	// In - where to look for the value
	In string

	// Name - matcher key
	Name string

	// Value - matcher value
	Value string
}
