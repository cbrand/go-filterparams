package filterparams

import (
	. "gopkg.in/check.v1"
)

var _ = Suite(&BuilderTest{})

type BuilderTest struct {
	builder *QueryBuilder
}

func (t *BuilderTest) SetUpTest(c *C) {
	t.builder = NewBuilder()
}

func (t *BuilderTest) TestEnableFilter(c *C) {
	t.builder.EnableFilter(FilterEq)
	c.Assert(t.builder.HasFilter(FilterEq.Identification), Equals, true)
}

func (t *BuilderTest) TestHasFilterNeg(c *C) {
	t.builder.EnableFilter(FilterGte)
	c.Assert(t.builder.HasFilter(FilterEq.Identification), Equals, false)
}

func (t *BuilderTest) TestGetFilter(c *C) {
	t.builder.EnableFilter(FilterEq)
	filter, err := t.builder.GetFilter(FilterEq.Identification)
	c.Assert(err, IsNil)
	c.Assert(filter, Equals, FilterEq)
}

func (t *BuilderTest) TestGetFilterNegative(c *C) {
	t.builder.EnableFilter(FilterLike)
	filter, err := t.builder.GetFilter(FilterEq.Identification)
	c.Assert(err, NotNil)
	c.Assert(filter, IsNil)
}
