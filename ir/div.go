package ir

import (
	"go/ast"
	"strings"

	"github.com/JackDalberg/SeaOfNodes/ir/types"
)

type DivNode struct {
	expr ast.Expr
	binaryNode
}

func NewDivNode(expr ast.Expr, lhs, rhs Node) *DivNode {
	return initBinaryNode(&DivNode{expr: expr}, lhs, rhs)
}

func (d *DivNode) GraphicLabel() string {
	return "/"
}

func (d *DivNode) label() string {
	return "Div"
}

func (d *DivNode) compute() (types.Type, error) {
	lType, ok := d.Lhs().base().typ.(*types.IntType)
	if !ok {
		return types.BottomType, nil
	}
	rType, ok := d.Rhs().base().typ.(*types.IntType)
	if !ok {
		return types.BottomType, nil
	}

	if lType.Constant() && rType.Constant() {
		if rType.Value == 0 {
			return nil, computeError(d.expr, "divide by zero")
		}
		return types.NewIntType(lType.Value / rType.Value), nil
	}
	return types.BottomType, nil
}

func (d *DivNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("(")
	toString(d.Lhs(), sb)
	sb.WriteString("/")
	toString(d.Rhs(), sb)
	sb.WriteString(")")
}
