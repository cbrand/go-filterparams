package filterparams

// Param is the absolute base of all statements in the entry.
type Param struct {

}

// LeftRight is a base struct for entries which have two entries.
type LeftRight struct {
	Param
	Left  *Param
	Right *Param
}

// And is an end statement for the entry.
type And struct {
	LeftRight
}

// Or represents an OR construct for the given entry.
type Or struct {
	LeftRight
}

// Argument is one of the
type Argument struct {
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

// Negate represents a negation of the entry.
type Negate struct {
	Param
	Statement *Param
}
