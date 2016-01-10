package filterparams

import (
	"regexp"

	"net/url"

	"github.com/cbrand/go-filterparams/definition"
)



var paramFilter = regexp.MustCompile("([a-zA-Z1-9]+)\\[([a-zA-Z1-9]+)\\]")

// Query can be used to parse query values.
type Query struct {
	filters    []*definition.Filter
}

func (q *Query) parseFilterArguments(values *url.Values) *ValueFilterArguments {
	arguments := NewValueFilterArgument()

	for param, valueList := range *values {
		for _, value := range valueList {
			matches := paramFilter.FindStringSubmatch(param)
			if matches == nil {
				continue
			}
			category, innerField := matches[0], matches[1]
			if category != "filter" {
				continue
			}

			matches = paramFilter.FindStringSubmatch(innerField)
			if matches != nil {
				innerCategory, innerField := matches[0], matches[1]
				if innerCategory == "param" {
					arguments.SetArgument(innerField, value)
				}
			} else if(innerField == "binding") {
				arguments.SetQueryBinding(value)
			} else if(innerField == "order") {
				arguments.AddOrder(value)
			}
		}
	}

	return arguments
}

/*func (q *Query) Parse(values *url.Values) *definition.Param {
	arguments := q.parseFilterArguments(values)
	if arguments.HasQueryBinding() {

	}
}

func (q *Query) parseFilter(field, value string) *definition.Parameter {

}*/

// newQuery uses the QueryBuilder to create a new Query entry.
func newQuery(allowedFilters []*definition.Filter) *Query {
	return &Query{
		filters: allowedFilters,
	}
}
