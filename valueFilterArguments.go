package filterparams

import (
	"fmt"
	"strings"
	"regexp"

	"github.com/cbrand/go-filterparams/binding"
	"github.com/cbrand/go-filterparams/definition"
)

// orderMatcher is used to verify
var orderMatcher = regexp.MustCompile("(?:(asc|desc)\\(([a-zA-Z0-9]*)\\)|([a-zA-Z0-9]*))")

// ParamNotFoundError represents a parameter which is specified
// in the query.
type ParamNotFoundError struct {
	ParamName string
}

// Error returns the string representation of the given error.
func (f *ParamNotFoundError) Error() string {
	return fmt.Sprintf("Parameter \"%s\" missing", f.ParamName)
}

// NewFilterParamNotFoundError creates a new error with the given parameter name.
func NewFilterParamNotFoundError(paramName string) *ParamNotFoundError {
	return &ParamNotFoundError{
		ParamName: paramName,
	}
}

// ValueFilterArguments is an internal used struct to adjust the filter arguments
// entry.
type ValueFilterArguments struct {
	arguments    map[string]*definition.Parameter
	queryBinding string
	orders       []string
}

// GetArgument returns the value of the arugment with the given name. Returns an empty
// string if the argument is not present
func (v *ValueFilterArguments) GetArgument(key string) *definition.Parameter {
	return v.arguments[key]
}

// SetArgument sets the argument with the passed data.
func (v *ValueFilterArguments) SetArgument(key string, value *definition.Parameter) {
	v.arguments[key] = value
}

// DelArgument removes the argument with the given entry.
func (v *ValueFilterArguments) DelArgument(key string) {
	delete(v.arguments, key)
}

// SetQueryBinding sets the string representation of the binding
// for the given query.
func (v *ValueFilterArguments) SetQueryBinding(binding string) {
	v.queryBinding = binding
}

// GetQueryBinding returns the set queryBinding of this configuartion or an
// empty string if none has been set.
func (v *ValueFilterArguments) GetQueryBinding() string {
	return v.queryBinding
}

// HasQueryBinding returns if an query binding has been set.
func (v *ValueFilterArguments) HasQueryBinding() bool {
	return len(v.GetQueryBinding()) > 0
}

// AddOrder adds an order entry to the system.
func (v *ValueFilterArguments) AddOrder(order string) {
	v.orders = append(v.orders, order)
}

// GetOrders returns the configured ordered parameters.
func (v *ValueFilterArguments) GetOrders() []string {
	return v.orders
}

// ParsedBinding parses the order string and returns the parsed result.
func (v *ValueFilterArguments) ParsedBinding() (interface{}, error) {
	data, err := binding.ParseString(v.queryBinding)
	if err != nil {
		return nil, err
	}
	parameters := data.(definition.ParameterHaver).GetParameters()
	for _, parameter := range parameters {
		argument := v.GetArgument(parameter.Identification)
		if argument == nil {
			return nil, NewFilterParamNotFoundError(parameter.Identification)
		}

		parameter.Name = argument.Name
		parameter.Value = argument.Value
		parameter.Filter = argument.Filter
	}

	return data, nil
}

// ApplyOrders takes the configured orders and returns the configured
// order objects.
func (v *ValueFilterArguments) ApplyOrders() []*definition.Order {
	orders := []*definition.Order{}
	for _, orderString := range v.GetOrders() {
		matches := orderMatcher.FindStringSubmatch(orderString)
		if len(matches) == 0 {
			continue
		}
		ascDesc := "asc"
		name := matches[3]
		if len(matches[3]) == 0 {
			ascDesc = matches[1]
			name = matches[2]
		}
		orders = append(orders, definition.NewOrder(name, ascDesc))
	}
	return orders
}

// ConstructDefaultQueryBinding creates a query binding where all parameters
// are connected with an AND-Statement.
func (v *ValueFilterArguments) ConstructDefaultQueryBinding() string {
	args := v.arguments
	data := make([]string, len(args))
	index := 0
	for key := range args {
		data[index] = key
		index++
	}
	return strings.Join(data, "&")
}

// NewValueFilterArgument initializes the ValueFilterArguments struct
// which is then used to store the parsed data.
func NewValueFilterArgument() *ValueFilterArguments {
	return &ValueFilterArguments{
		arguments: map[string]*definition.Parameter{},
		orders: []string{},
	}
}
