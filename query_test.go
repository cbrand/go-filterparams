package filterparams

import (
	"fmt"
	"net/url"

	. "gopkg.in/check.v1"

	"github.com/cbrand/go-filterparams/definition"
)

var _ = Suite(&QueryTest{})

type QueryTest struct {
	builder *QueryBuilder
	data    *url.Values
}

func (t *QueryTest) SetUpTest(c *C) {
	t.builder = NewBuilder()
	t.builder.EnableFilter(definition.FilterEq)
	t.builder.EnableFilter(definition.FilterLike)
	t.data = &url.Values{}
}

func (t *QueryTest) TestGetFilterArguments(c *C) {
	urlTemplate := "http://myurl.com?filter[param][reference][eq][introducer_name]=%s&filter[param][references][eq][agreement_number]=%s&filter[binding]=%s"
	urlString := fmt.Sprintf(urlTemplate, url.QueryEscape("Broker 1"), url.QueryEscape("123456789"), url.QueryEscape("(introducer_name&agreement_number)"))
	query := t.builder.CreateQuery()

	toParseURL, err := url.Parse(urlString)
	c.Assert(err, IsNil)
	values := toParseURL.Query()
	queryData, err := query.Parse(&values)
	c.Assert(err, IsNil)
	filter := queryData.GetFilter()
	and, ok := filter.(*definition.And)
	c.Assert(ok, Equals, true)
	parameter, ok := and.Left.(*definition.Parameter)
	c.Assert(ok, Equals, true)

	c.Assert(parameter.Identification, Equals, "introducer_name")
	c.Assert(parameter.Value, Equals, "Broker 1")

	parameter, ok = and.Right.(*definition.Parameter)
	c.Assert(ok, Equals, true)
	c.Assert(parameter.Identification, Equals, "agreement_number")
	c.Assert(parameter.Value, Equals, "123456789")
}

func (t *QueryTest) run(c *C) *QueryData {
	queryData, err := t.builder.CreateQuery().Parse(t.data)
	c.Assert(err, IsNil)
	return queryData
}

func (t *QueryTest) expectError(c *C) {
	_, err := t.builder.CreateQuery().Parse(t.data)
	c.Assert(err, NotNil)
}

func (t *QueryTest) constructFilterParamKey(name, operation, alias string) string {
	if len(name) <= 0 {
		panic("Name must be set")
	}
	param := fmt.Sprintf("filter[param][%s]", name)
	if len(alias) > 0 && len(operation) == 0 {
		operation = "eq"
	}

	for _, toAddItem := range []string{operation, alias} {
		if len(toAddItem) > 0 {
			param = fmt.Sprintf("%s[%s]", param, toAddItem)
		}
	}
	return param
}

func (t *QueryTest) addAliasedFilterParam(name, operation, alias, value string) {
	param := t.constructFilterParamKey(name, operation, alias)
	t.data.Set(param, value)
}

func (t *QueryTest) addFilterParam(name, operation, value string) {
	t.addAliasedFilterParam(name, operation, "", value)
}

func (t *QueryTest) addDateParam() {
	t.addFilterParam("date", "eq", "2015-01-01")
	t.data.Set("filter[param][date][eq]", "2015-01-01")
}

func (t *QueryTest) addNameParam() {
	t.addFilterParam("name", "like", "%helloWorld%")
}

func (t *QueryTest) addNameAndDateParam() {
	t.addDateParam()
	t.addNameParam()
}

func (t *QueryTest) TestQuery(c *C) {
	t.addDateParam()
	queryData := t.run(c)
	filter := queryData.GetFilter()
	param, ok := filter.(*definition.Parameter)
	c.Assert(ok, Equals, true)
	c.Assert(param.Filter, Equals, definition.FilterEq)
	c.Assert(param.Name, Equals, "date")
	c.Assert(param.Identification, Equals, "date")
	c.Assert(param.Value, Equals, "2015-01-01")
}

func (t *QueryTest) TestQueryUnsupportedOperation(c *C) {
	t.data.Set("filter[param][date][notSupported]", "2015-01-01")
	t.expectError(c)
}

func (t *QueryTest) TestQueryWithImplicitBinding(c *C) {
	t.addNameAndDateParam()
	queryData := t.run(c)
	filter := queryData.GetFilter()
	param, ok := filter.(*definition.And)
	c.Assert(ok, Equals, true)
	leftData, ok := param.Left.(*definition.Parameter)
	c.Assert(ok, Equals, true)
	rightData, ok := param.Right.(*definition.Parameter)
	c.Assert(ok, Equals, true)
	dateData := leftData
	nameData := rightData
	if dateData.Name != "date" {
		nameData, dateData = dateData, nameData
	}

	c.Assert(dateData.Name, Equals, "date")
	c.Assert(dateData.Filter, Equals, definition.FilterEq)
	c.Assert(dateData.Value, Equals, "2015-01-01")
	c.Assert(nameData.Name, Equals, "name")
	c.Assert(nameData.Filter, Equals, definition.FilterLike)
	c.Assert(nameData.Value, Equals, "%helloWorld%")
}

func (t *QueryTest) TestQueryWithExplicitBinding(c *C) {
	t.addNameAndDateParam()
	t.data.Set("filter[binding]", "name|!date")
	queryData := t.run(c)
	filter := queryData.GetFilter()
	or, ok := filter.(*definition.Or)
	c.Assert(ok, Equals, true)
	c.Assert(ok, Equals, true)

	leftData, ok := or.Left.(*definition.Parameter)
	c.Assert(ok, Equals, true)
	rightNegation, ok := or.Right.(*definition.Negate)
	c.Assert(ok, Equals, true)
	rightData, ok := rightNegation.Negated.(*definition.Parameter)
	c.Assert(ok, Equals, true)
	c.Assert(leftData.Name, Equals, "name")
	c.Assert(rightData.Name, Equals, "date")
}

func (t *QueryTest) TestQueryWithSetAlias(c *C) {
	t.addAliasedFilterParam("name", "eq", "aliasedName", "smith")
	queryData := t.run(c)
	filter := queryData.GetFilter()
	param, ok := filter.(*definition.Parameter)
	c.Assert(ok, Equals, true)
	c.Assert(param.Identification, Equals, "aliasedName")
	c.Assert(param.Name, Equals, "name")
}

func (t *QueryTest) TestQueryWithSetAliasAndBinding(c *C) {
	t.addAliasedFilterParam("name", "eq", "aliasedName", "smith")
	t.addAliasedFilterParam("status", "like", "aliasedStatus", "%single%")
	t.data.Set("filter[binding]", "(aliasedName|aliasedStatus)|!aliasedName")
	queryData := t.run(c)
	filter := queryData.GetFilter()
	or, ok := filter.(*definition.Or)
	c.Assert(ok, Equals, true)
	rightNegation := or.Right.(*definition.Negate)
	parameter, ok := rightNegation.Negated.(*definition.Parameter)
	c.Assert(ok, Equals, true)
	c.Assert(parameter.Identification, Equals, "aliasedName")
	c.Assert(parameter.Name, Equals, "name")

	leftOr, ok := or.Left.(*definition.Or)
	c.Assert(ok, Equals, true)
	leftOrAliasedName, ok := leftOr.Left.(*definition.Parameter)
	c.Assert(ok, Equals, true)
	c.Assert(leftOrAliasedName.Identification, Equals, "aliasedName")
	c.Assert(leftOrAliasedName.Name, Equals, "name")

	leftOrAliasedStatus, ok := leftOr.Right.(*definition.Parameter)
	c.Assert(ok, Equals, true)
	c.Assert(leftOrAliasedStatus.Identification, Equals, "aliasedStatus")
	c.Assert(leftOrAliasedStatus.Name, Equals, "status")
}

func (t *QueryTest) TestQueryEmptyOrder(c *C) {
	t.addNameParam()
	queryData := t.run(c)

	c.Assert(len(queryData.GetOrders()), Equals, 0)
}

func (t *QueryTest) addOrder(order string) {
	t.data.Add("filter[order]", order)
}

func (t *QueryTest) TestQueryWithOrder(c *C) {
	t.addNameParam()
	t.addOrder("name")
	t.addOrder("desc(date)")
	queryData := t.run(c)
	orders := queryData.GetOrders()
	c.Assert(len(orders), Equals, 2)

	orderNames := make([]string, len(orders))
	orderDesc := make([]bool, len(orders))
	for index, order := range orders {
		orderNames[index] = order.GetOrderBy()
		orderDesc[index] = order.OrderDesc()
	}

	c.Assert(orderNames, DeepEquals, []string{"name", "date"})
	c.Assert(orderDesc, DeepEquals, []bool{false, true})
}

func (t *QueryTest) TestIgnoreOtherParams(c *C) {
	t.addNameParam()
	t.data.Set("filter[unrecognized]", "hallo")
	t.data.Set("nonFilter", "test")
	t.data.Set("other[category][differ]", "egh")
	queryData := t.run(c)
	_, ok := queryData.GetFilter().(*definition.Parameter)
	c.Assert(ok, Equals, true)
}
