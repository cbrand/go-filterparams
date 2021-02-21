package filterparams

import (
	"regexp"

	"net/url"

	"github.com/cbrand/go-filterparams/definition"
)

const defaultOperation = "eq"

var paramFilter = regexp.MustCompile("([a-zA-Z1-9_\\-]+)\\[([a-zA-Z1-9_\\-]+)\\](.*)")
var fieldFilter = regexp.MustCompile("\\[([a-zA-Z1-9_\\-]+)\\](.*)")

// Query can be used to parse query values.
type Query struct {
	filters []*definition.Filter
	defaultOperation string
}

// parseFilterArguments takes the filter arugments and parses the data.
func (q *Query) parseFilterArguments(values *url.Values) (*ValueFilterArguments, error) {
	arguments := NewValueFilterArgument()

	for param, valueList := range *values {
		matches := paramFilter.FindStringSubmatch(param)
		if len(matches) == 0 {
			continue
		}
		category, innerCategory, remaining := matches[1], matches[2], matches[3]
		if category != "filter" {
			continue
		}

		for _, value := range valueList {
			matches = fieldFilter.FindStringSubmatch(remaining)

			if innerCategory == "param" {
				parameter, err := q.parseFilterParam(matches[1], matches[2], value)
				if err != nil {
					return nil, err
				}
				arguments.SetArgument(parameter.Identification, parameter)
			} else if (innerCategory == "binding") {
				arguments.SetQueryBinding(value)
			} else if (innerCategory == "order") {
				arguments.AddOrder(value)
			}
		}
	}

	return arguments, nil
}

// parseFilterParam takes the basic configuration and generates a filter parameter.
func (q *Query) parseFilterParam(paramName, remainingKeyData, value string) (*definition.Parameter, error) {
	remainingDataMatches := fieldFilter.FindStringSubmatch(remainingKeyData)
	operation := defaultOperation
	possibleRemainingAlias := ""
	if remainingDataMatches != nil {
		operation = remainingDataMatches[1]
		possibleRemainingAlias = remainingDataMatches[2]
	}
	remainingAliasMatches := fieldFilter.FindStringSubmatch(possibleRemainingAlias)
	alias := paramName
	if remainingAliasMatches != nil {
		alias = remainingAliasMatches[1]
	}
	parameter := definition.NewParameter(alias)
	parameter.Name = paramName
	parameter.Filter = q.getFilter(operation)
	parameter.Value = value
	if parameter.Filter == nil {
		return nil, NewUnsupportedOperation(operation)
	}

	return parameter, nil
}

// getFilter returns the registered filter with the given name.
func (q *Query) getFilter(name string) *definition.Filter {
	for _, filter := range q.filters {
		if filter.Identification == name {
			return filter
		}
	}
	return nil
}

// GetDefaultOperation returns the default operation as the
// entry.
func (q *Query) GetDefaultOperation() string {
	operation := defaultOperation
	if len(q.defaultOperation) > 0 {
		operation = q.defaultOperation
	}
	return operation
}

// setDefaultOperation is used by the builder to be able to
// configure a default operation.
func (q *Query) setDefaultOperation(operation string) {
	q.defaultOperation = operation
}

// Parse takes the given values and returns the parsed data which is provided
// by the Go struct.
func (q *Query) Parse(values *url.Values) (*QueryData, error) {
	arguments, err := q.parseFilterArguments(values)

	if err != nil {
		return nil, err
	}

	if !arguments.HasQueryBinding() {
		arguments.SetQueryBinding(arguments.ConstructDefaultQueryBinding())
	}
	binding, err := arguments.ParsedBinding()
	if err != nil {
		return nil, err
	}

	orders := arguments.ApplyOrders()

	return NewQueryData(binding, orders), nil
}

// newQuery uses the QueryBuilder to create a new Query entry.
func newQuery(allowedFilters []*definition.Filter) *Query {
	return &Query{
		filters: allowedFilters,
	}
}
