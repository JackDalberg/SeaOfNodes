package ir

import (
	"go/types"
	"strconv"
	"strings"
)

type ConstantNode struct {
	baseNode
	Value int
}

func NewConstantNode(typ types.Type) *ConstantNode {
	n := initBaseNode(&ConstantNode{}, StartNode)
	n.typ = typ
	return n
}

func (c *ConstantNode) IsControl() bool {
	return false
}

func (c *ConstantNode) compute() (types.Type, error) {
	return c.typ, nil
}

func (c *ConstantNode) label() string {
	return "#" + strconv.Itoa(c.value())
}

func (c *ConstantNode) GraphicLabel() string {
	return c.label()
}

func (c *ConstantNode) toStringInternal(sb *strings.Builder) {
	c.typ.ToString(sb)
}

func (c *ConstantNode) value() int {
	return c.typ.(*types.IntValue).Value
}
