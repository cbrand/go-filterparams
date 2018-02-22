package binding

import (
	"github.com/cbrand/go-filterparams/definition"

	. "gopkg.in/check.v1"
)

var _ = Suite(&ParseTest{})

type ParseTest struct{}

func (t *ParseTest) expectParsingError(c *C, data string) {
	_, err := ParseString(data)
	c.Assert(err, NotNil)
}

func (t *ParseTest) parse(c *C, data string) interface{} {
	item, err := ParseString(data)
	c.Assert(err, IsNil)
	return item
}

func (t *ParseTest) TestParseTerm(c *C) {
	data := t.parse(c, "as")
	param, ok := data.(*definition.Parameter)
	c.Assert(ok, Equals, true)
	c.Assert(param.Identification, Equals, "as")
}

func (t *ParseTest) TestParseNegativeTerm(c *C) {
	data := t.parse(c, "! as")
	param, ok := data.(*definition.Negate)
	c.Assert(ok, Equals, true)
	c.Assert(param.Negated.(*definition.Parameter).Identification, Equals, "as")
}

func (t *ParseTest) TestParseBracket(c *C) {
	data := t.parse(c, "( as )")
	param, ok := data.(*definition.Parameter)
	c.Assert(ok, Equals, true)
	c.Assert(param.Identification, Equals, "as")
}

func (t *ParseTest) TestParsePartialBracketError(c *C) {
	t.expectParsingError(c, "( as")
}

func (t *ParseTest) TestParseAnd(c *C) {
	data := t.parse(c, "left &right")
	and, ok := data.(*definition.And)
	c.Assert(ok, Equals, true)
	left, right := and.Left, and.Right
	c.Assert(left.(*definition.Parameter).Identification, Equals, "left")
	c.Assert(right.(*definition.Parameter).Identification, Equals, "right")
}

func (t *ParseTest) TestParseOr(c *C) {
	data := t.parse(c, "left| right")
	or, ok := data.(*definition.Or)
	c.Assert(ok, Equals, true)
	left, right := or.Left, or.Right
	c.Assert(left.(*definition.Parameter).Identification, Equals, "left")
	c.Assert(right.(*definition.Parameter).Identification, Equals, "right")
}

func (t *ParseTest) TestParsePartialAnd(c *C) {
	t.expectParsingError(c, "a &")
}

func (t *ParseTest) TestParsePartialOr(c *C) {
	t.expectParsingError(c, "a |")
}

func (t *ParseTest) TestParseParamSpace(c *C) {
	data := t.parse(c, " as       \t\n   ")
	param, ok := data.(*definition.Parameter)
	c.Assert(ok, Equals, true)
	c.Assert(param.Identification, Equals, "as")
}

func (t *ParseTest) TestParseAndOr(c *C) {
	data := t.parse(c, "left & (middle | data) | !right")
	outerOr, ok := data.(*definition.Or)
	c.Assert(ok, Equals, true)
	rightNot, ok := outerOr.Right.(*definition.Negate)
	c.Assert(ok, Equals, true)
	c.Assert(rightNot.Negated.(*definition.Parameter).Identification, Equals, "right")

	leftAnd, ok := outerOr.Left.(*definition.And)
	c.Assert(ok, Equals, true)

	leftArgument, ok := leftAnd.Left.(*definition.Parameter)
	c.Assert(leftArgument.Identification, Equals, "left")

	rightInnerOr, ok := leftAnd.Right.(*definition.Or)
	c.Assert(ok, Equals, true)
	rightInnerLeft, ok := rightInnerOr.Left.(*definition.Parameter)
	c.Assert(ok, Equals, true)
	c.Assert(rightInnerLeft.Identification, Equals, "middle")
	rightInnerRight, ok := rightInnerOr.Right.(*definition.Parameter)
	c.Assert(ok, Equals, true)
	c.Assert(rightInnerRight.Identification, Equals, "data")
}

func (t *ParseTest) TestParseBrackets(c *C) {
	data := t.parse(c, "(left1 | left2) & (right1 | right2)")
	outerAnd, ok := data.(*definition.And)
	c.Assert(ok, Equals, true)
	_, ok = outerAnd.Left.(*definition.Or)
	c.Assert(ok, Equals, true)
	_, ok = outerAnd.Right.(*definition.Or)
	c.Assert(ok, Equals, true)
}

func (t *ParseTest) TestParseOneBracket(c *C) {
	data := t.parse(c, "(left1 | left2) & right1 | right2")
	outerAnd, ok := data.(*definition.Or)
	c.Assert(ok, Equals, true)
	_, ok = outerAnd.Right.(*definition.Parameter)
	c.Assert(ok, Equals, true)
	leftAnd, ok := outerAnd.Left.(*definition.And)
	c.Assert(ok, Equals, true)
	_, ok = leftAnd.Right.(*definition.Parameter)
	c.Assert(ok, Equals, true)
	_, ok = leftAnd.Left.(*definition.Or)
	c.Assert(ok, Equals, true)
}

func (t *ParseTest) TestParseNegation(c *C) {
	data := t.parse(c, "!item1 & (item2 | item3)")
	outerAnd, ok := data.(*definition.And)
	c.Assert(ok, Equals, true)
	_, ok = outerAnd.Left.(*definition.Negate)
	c.Assert(ok, Equals, true)
	_, ok = outerAnd.Right.(*definition.Or)
	c.Assert(ok, Equals, true)
}

func (t *ParseTest) TestParseNegationAnd(c *C) {
	data := t.parse(c, "!(item1 & item2)")
	outerNegation, ok := data.(*definition.Negate)
	c.Assert(ok, Equals, true)
	_, ok = outerNegation.Negated.(*definition.And)
	c.Assert(ok, Equals, true)
}

func (t *ParseTest) TestParseNegationOr(c *C) {
	data := t.parse(c, "!(item1 | item2)")
	outerNegation, ok := data.(*definition.Negate)
	c.Assert(ok, Equals, true)
	_, ok = outerNegation.Negated.(*definition.Or)
	c.Assert(ok, Equals, true)
}
