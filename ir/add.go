package ir

import "strings"

type AddNode struct {
	binaryNode
}

func NewAddNode(lhs, rhs Node) *AddNode {
	return initBinaryNode(&AddNode{}, lhs, rhs)
}

func (a *AddNode) GraphicLabel() string {
	return "+"
}

func (a *AddNode) label() string {
	return "Add"
}

func (a *AddNode) compute() (types.Type, error) {
	lType, ok := a.Lhs().base().typ.(*types.IntType)
	if !ok {
		return types.Bottomtype, nil
	}
	rType, ok := a.Rhs().base().typ.(*types.IntType)
	if !ok {
		return types.Bottomtype, nil
	}

	if lType.Constant() && rType.Constant() {
		return types.NewIntType(lType.value + rType.Value), nil
	}
	return types.BottomType, nil
}

func (a *AddNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("(")
	toString(a.Lhs(), sb)
	sb.WriteString("+")
	toString(a.Rhs(), sb)
	sb.WriteString(")")
}
