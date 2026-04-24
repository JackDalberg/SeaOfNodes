package ir

import "strings"

type MulNode struct {
	binaryNode
}

func NewMulNode(lhs, rhs Node) *MulNode {
	return initBinaryNode(&MulNode{}, lhs, rhs)
}

func (m *MulNode) GraphicLabel() string {
	return "*"
}

func (m *MulNode) label() string {
	return "Mul"
}

func (m *MulNode) compute() (types.Type, error) {
	lType, ok := m.Lhs().base().typ.(*types.IntType)
	if !ok {
		return types.Bottomtype, nil
	}
	rType, ok := m.Rhs().base().typ.(*types.IntType)
	if !ok {
		return types.Bottomtype, nil
	}

	if lType.Constant() && rType.Constant() {
		return types.NewIntType(lType.value * rType.Value), nil
	}
	return types.BottomType, nil
}

func (m *MulNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("(")
	toString(m.Lhs(), sb)
	sb.WriteString("*")
	toString(m.Rhs(), sb)
	sb.WriteString(")")
}
