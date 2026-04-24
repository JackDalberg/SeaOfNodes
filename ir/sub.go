package ir

import "strings"

type SubNode struct {
	binaryNode
}

func NewSubNode(lhs, rhs Node) *SubNode {
	return initBinaryNode(&SubNode{}, lhs, rhs)
}

func (s *SubNode) GraphicLabel() string {
	return "-"
}

func (s *SubNode) label() string {
	return "Sub"
}

func (s *SubNode) compute() (types.Type, error) {
	lType, ok := s.Lhs().base().typ.(*types.IntType)
	if !ok {
		return types.Bottomtype, nil
	}
	rType, ok := s.Rhs().base().typ.(*types.IntType)
	if !ok {
		return types.Bottomtype, nil
	}

	if lType.Constant() && rType.Constant() {
		return types.NewIntType(lType.value - rType.Value), nil
	}
	return types.BottomType, nil
}

func (s *SubNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("(")
	toString(s.Lhs(), sb)
	sb.WriteString("-")
	toString(s.Rhs(), sb)
	sb.WriteString(")")
}
