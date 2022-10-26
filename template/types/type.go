package types

var types map[string]Type

type Type interface {
	Value() string
}
