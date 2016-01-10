package filterparams

// Query can be used to parse query values.
type Query struct {
	filters []*Filter
}

// newQuery uses the QueryBuilder to create a new Query entry.
func newQuery(allowedFilters []*Filter) *Query {
	return &Query{
		filters: allowedFilters,
	}
}
