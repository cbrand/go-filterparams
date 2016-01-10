package definition

// Param is the absolute base of all statements in the entry.
type Param struct {

}

// LeftRight is a base struct for entries which have two entries.
type LeftRight struct {
	Param
	Left  interface{}
	Right interface{}
}

// And is an end statement for the entry.
type And struct {
	LeftRight
}

// NewAnd returns a new empty And struct.
func NewAnd() *And {
	return &And{}
}

// Or represents an OR construct for the given entry.
type Or struct {
	LeftRight
}

// NewOr returns an empty or struct.
func NewOr() *Or {
	return &Or{}
}

// Parameter is one of the arguments which has a name, value and a filter
// which should be applied to the system.
type Parameter struct {
	Param
	// Identification is a name which could be set.
	Identification string
	// Name is the name of the field.
	Name           string
	// Filter is the filter entry for the given entry.
	Filter         *Filter
	// Value is the value which the entry should be filtered by.
	Value          interface{}
}

// NewParameter returns a new parameter initialized with the given
// identification.
func NewParameter(identification string) *Parameter {
	return &Parameter{
		Identification: identification,
	}
}

// Negate represents a negation of the entry.
type Negate struct {
	Param
	Statement *Parameter
}

// NewNegate returns the negation of the given element.
func NewNegate(statement *Parameter) *Negate {
	return &Negate{Statement: statement}
}
