package filterparams

import (
	"github.com/cbrand/go-filterparams/definition"

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
	t.builder.EnableFilter(definition.FilterEq)
	c.Assert(t.builder.HasFilter(definition.FilterEq.Identification), Equals, true)
}

func (t *BuilderTest) TestHasFilterNeg(c *C) {
	t.builder.EnableFilter(definition.FilterGte)
	c.Assert(t.builder.HasFilter(definition.FilterEq.Identification), Equals, false)
}

func (t *BuilderTest) TestGetFilter(c *C) {
	t.builder.EnableFilter(definition.FilterEq)
	filter, err := t.builder.GetFilter(definition.FilterEq.Identification)
	c.Assert(err, IsNil)
	c.Assert(filter, Equals, definition.FilterEq)
}

func (t *BuilderTest) TestGetFilterNegative(c *C) {
	t.builder.EnableFilter(definition.FilterLike)
	filter, err := t.builder.GetFilter(definition.FilterEq.Identification)
	c.Assert(err, NotNil)
	c.Assert(filter, IsNil)
}
