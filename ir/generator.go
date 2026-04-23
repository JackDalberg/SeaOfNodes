package ir

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

type ASTError struct {
	err error
	Pos token.Pos
}

func (a *ASTError) Error() string {
	return a.err.Error()
}

func astError(pos token.Pos, n ast.Node) *ASTError {
	err := fmt.Errorf("Unsupported AST: %v", n)
	return &ASTError{
		err: err,
		Pos: pos,
	}
}

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(n ast.Node) (*ReturnNode, error) {
	var retNode *ReturnNode
	var err error
	ast.Inspect(n, func(n ast.Node) bool {
		if err != nil {
			return false
		}

		if ret, ok := n.(*ast.ReturnStmt); ok {
			retNode, err = g.generateReturn(ret)
			if err != nil {
				return false
			}
		}
		return true
	})
	return retNode, err
}

func (g *Generator) generateReturn(r *ast.ReturnStmt) (*ReturnNode, error) {
	expr, err := g.generateExpr(r.Results[0])
	if err != nil {
		return nil, err
	}
	return NewReturnNode(StartNode, expr), nil
}

func (g *Generator) generateExpr(e ast.Expr) (Node, error) {
	switch t := e.(type) {
	case *ast.BasicLit:
		num, err := strconv.Atoi(t.Value)
		if err != nil {
			return nil, err
		}
		return NewConstantNode(num), nil

	default:
		return nil, astError(e.Pos(), e)
	}
}
