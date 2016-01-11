package filterparams

import "github.com/cbrand/go-filterparams/definition"

// QueryData represents the completely parsed Query and is returned
// after successful parsing is done.
type QueryData struct {
	// Filter represents the structure which should be used to filter
	// the set data.
	filter interface{}

	// Order is a sorted slice of order statements which should be
	// sorted by.
	order  []*definition.Order
}

// GetFilter returns the parsed filter of the QueryData.
func (q *QueryData) GetFilter() interface{} {
	return q.filter
}

// GetOrders returns a list of all items which the query should
// be ordered by.
func (q *QueryData) GetOrders() []*definition.Order {
	return q.order
}

// NewQueryData initializes a new QueryData struct with all filter information
// of the parameter.
func NewQueryData(filter interface{}, order []*definition.Order) *QueryData {
	return &QueryData{
		filter: filter,
		order: order,
	}
}
