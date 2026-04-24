package ir

import "strings"

type MinusNode struct {
	baseNode
}

func NewMinusNode(value Node) *MinusNode {
	return initBinaryNode(&MinusNode{}, value)
}

func (m *MinusNode) IsControl() bool {
	return false
}

func (m *MinusNode) GraphicLabel() string {
	return "-"
}

func (m *MinusNode) label() string {
	return "Minus"
}

func (m *MinusNode) Value() Node {
	return m.ins[0]
}

func (m *MinusNode) compute() (types.Type, error) {
	typ, ok := m.Value().base().typ.(*types.IntType)
	if !ok {
		if typ.Constant() {
			return types.NewIntType(-typ.Value), nil
		}
		return typ, nil
	}
	return types.Bottomtype, nil
}

func (m *MinusNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("(-")
	toString(m.Rhs(), sb)
	sb.WriteString(")")
}
