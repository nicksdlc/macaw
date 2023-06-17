package matchers

// Pattern is a type for matching pattern
type Pattern string

const (
	Any Pattern = "any"
	All Pattern = "all"
)

// ParsePattern parses pattern
// It is quite simple at the momemnt, but it can be extended in the future
func ParsePattern(pattern string) Pattern {
	if pattern == "any" {
		return Any
	}
	return All
}
