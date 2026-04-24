package ir

import (
	"go/ast"
	"strings"
)

type DivNode struct {
	expr ast.Expr
	baseNode
}

func NewDivNode(expr ast.Expr, lhs, rhs Node) *DivNode {
	return initBinaryNode(&DivNode{expr: expr}, lhs, rhs)
}

func (a *DivNode) GraphicLabel() string {
	return "/"
}

func (a *DivNode) label() string {
	return "Div"
}

func (a *DivNode) compute() (types.Type, error) {
	lType, ok := a.Lhs().base().typ.(*types.IntType)
	if !ok {
		return types.Bottomtype, nil
	}
	rType, ok := a.Rhs().base().typ.(*types.IntType)
	if !ok {
		return types.Bottomtype, nil
	}

	if lType.Constant() && rType.Constant() {
		if rType.Value == 0 {
			return nil, computeError(d.expr, "divide by zero")
		}
		return types.NewIntType(lType.value / rType.Value), nil
	}
	return types.BottomType, nil
}

func (a *DivNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("(")
	toString(a.Lhs(), sb)
	sb.WriteString("/")
	toString(a.Rhs(), sb)
	sb.WriteString(")")
}
