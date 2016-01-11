package definition

import (
		. "gopkg.in/check.v1"
)

var _ = Suite(&ParamTest{})

type ParamTest struct {
	Parameter ParameterHaver
}

func (t *ParamTest) SetUpTest(c *C) {
	data := NewAnd()
	data.Left = NewParameter("left")
	data.Right = NewNegate(NewParameter("right"))

	t.Parameter = data
}

func (t *ParamTest) TestGetParameters(c *C) {
	param := t.Parameter
	parameters := param.GetParameters()
	c.Assert(len(parameters), Equals, 2)
	c.Assert(parameters[0].Identification, Equals, "left")
	c.Assert(parameters[1].Identification, Equals, "right")
}
