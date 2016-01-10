package filterparams

import (
	"net/url"
)

// ValueFilterArguments is an internal used struct to adjust the filter arguments
// entry.
type ValueFilterArguments struct {
	arguments    *url.Values
	queryBinding string
	orders       []string
}

// GetArgument returns the value of the arugment with the given name. Returns an empty
// string if the argument is not present
func (v *ValueFilterArguments) GetArgument(key string) string {
	return v.arguments.Get(key)
}

// SetArgument sets the argument with the passed data.
func (v *ValueFilterArguments) SetArgument(key, value string) {
	v.arguments.Set(key, value)
}

// DelArgument removes the argument with the given entry.
func (v *ValueFilterArguments) DelArgument(key string) {
	v.arguments.Del(key)
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

func NewValueFilterArgument() *ValueFilterArguments {
	return &ValueFilterArguments{
		arguments: new(url.Values),
		orders: []string{},
	}
}
