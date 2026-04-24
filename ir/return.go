package ir

import (
	"strings"

	"github.com/JackDalberg/SeaOfNodes/ir/types"
)

type ReturnNode struct {
	baseNode
}

func NewReturnNode(control, data Node) *ReturnNode {
	return initBaseNode(&ReturnNode{}, control, data)
}

func (r *ReturnNode) IsControl() bool {
	return true
}

func (r *ReturnNode) Control() Node {
	return In(r, 0)
}

func (r *ReturnNode) Expr() Node {
	return In(r, 1)
}

func (r *ReturnNode) label() string {
	return "Return"
}

func (r *ReturnNode) GraphicLabel() string {
	return r.label()
}

func (r *ReturnNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("return ")
	toString(r.Expr(), sb)
	sb.WriteString(";")
}
func (r *ReturnNode) compute() (types.Type, error) {
	return types.BottomType, nil
}
