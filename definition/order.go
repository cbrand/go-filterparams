package definition

import "strings"

type Order struct {
	orderBy string
	orderDesc bool
}

// GetOrderBy returns the parameter name it should be ordered by.
func (o *Order) GetOrderBy() string {
	return o.orderBy
}

// OrderDesc returns if the sorting should be ordered by in descending order.
func (o *Order) OrderDesc() bool {
	return o.orderDesc
}

func newOrder(orderBy string) *Order {
	return &Order{
		orderBy: orderBy,
	}
}

// NewOrderAsc returns a new order object which values should be used to order
// by ascending value.
func NewOrderAsc(orderBy string) *Order {
	order := newOrder(orderBy)
	order.orderDesc = false
	return order
}

// NewOrderDesc returns a new order object which values should be used to order
// by descending order.
func NewOrderDesc(orderBy string) *Order {
	order := newOrder(orderBy)
	order.orderDesc = true
	return order
}

// NewOrder generates a new order instance with the given sortOrder. If the
// sortOrder is set to "desc" than it is marked to be sorted in descending order.
// In any other instance it is sorted in ascending order.
func NewOrder(orderBy string, sortOrder string) *Order {
	if strings.ToLower(sortOrder) == "desc" {
		return NewOrderDesc(orderBy)
	} else {
		return NewOrderAsc(orderBy)
	}
}
