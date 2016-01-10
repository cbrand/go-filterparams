package binding

// ParseString parses the given query order parameter and returns a parameter
// entry.
func ParseString(data string, opts ...Option) (interface{}, error) {
	return Parse("data", []byte(data), opts...)
}
