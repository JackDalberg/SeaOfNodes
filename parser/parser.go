package parser

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"
)

type SyntaxError struct {
	err    error
	Offset int
}

func (s *SyntaxError) Error() string {
	return s.err.Error()
}

func syntaxError(offset int, msgFormat string, args ...any) *SyntaxError {
	err := fmt.Errorf("Syntax error: "+msgFormat, args...)
	return &SyntaxError{
		err:    err,
		Offset: offset,
	}
}

type Parser struct {
	lexer   *lexer
	file    *token.File
	fileset *token.FileSet
	source  string
}

func NewParser(source string) *Parser {
	fileset := token.NewFileSet()
	return &Parser{
		lexer:   NewLexer(source),
		file:    fileset.AddFile("", 1, len(source)),
		fileset: fileset,
		source:  source,
	}
}

func (p *Parser) Parse() (*ast.ReturnStmt, error) {
	n, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	if b, offset, ok := p.lexer.ReadByte(); ok {
		return nil, syntaxError(offset, "unexpected %c", b)
	}
	return n.(*ast.ReturnStmt), nil
}

func (p *Parser) parseStatement() (ast.Node, error) {
	t, offset, err := p.lexer.ReadToken()
	if err != nil {
		return nil, err
	}

	switch t {
	case "return":
		return p.parseReturn(p.offsetToPos(offset))

	default:
		return nil, syntaxError(offset, "expected a statement, got=%s", t)
	}
}

func (p *Parser) parseReturn(pos token.Pos) (*ast.ReturnStmt, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	token, offset, ok := p.lexer.ReadByte()
	if !ok {
		return nil, syntaxError(offset, "expected ';' after expression")
	}
	if token != ';' {
		return nil, syntaxError(offset, "expected ';', got=%c", token)
	}

	return &ast.ReturnStmt{Return: p.offsetToPos(offset), Results: []ast.Expr{expr}}, nil
}

func (p *Parser) parseExpr() (ast.Expr, error) {
	return p.parseBinary()
}

func (p *Parser) parsePrimary() (*ast.BasicLit, error) {
	num, offset, err := p.lexer.ReadNumber()
	if err != nil {
		return nil, syntaxError(offset, err.Error())
	}
	blit := &ast.BasicLit{
		ValuePos: p.offsetToPos(offset),
		Kind:     token.INT,
		Value:    num,
	}
	return blit, nil
}

func (p *Parser) offsetToPos(offset int) token.Pos {
	return p.file.Pos(offset)
}

func (p *Parser) PosToOffset(pos token.Pos) int {
	return p.file.Offset(pos)
}

func opToToken(op byte) token.Token {
	switch op {
	case '+':
		return token.ADD
	case '-':
		return token.SUB
	case '*':
		return token.MUL
	case '/':
		return token.QUO
	default:
		return token.ILLEGAL
	}
}

func (p *Parser) string(n ast.Node) string {
	sb := &strings.Builder{}
	printer.Fprint(sb, p.fileset, n)
	return sb.String()
}

func (p *Parser) parseUnary() (ast.Expr, error) {
	op, opPos, hasOp := p.parseOp()
	if !hasOp {
		return p.parsePrimary()
	}

	value, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	return &ast.UnaryExpr{X: value, Op: op, OpPos: opPos}, nil
}

func (p *Parser) parseBinary() (ast.Expr, error) {
	lhs, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	return p.parseRhs(lhs)
}

func (p *Parser) parseRhs(lhs ast.Expr) (ast.Expr, error) {
	op, opPos, hasOp := p.parseOp()
	if !hasOp {
		return lhs, nil
	}

	rhs, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	return p.parseRhs(p.withPrecedence(lhs, op, opPos, rhs))
}

func (p *Parser) parseOp() (token.Token, token.Pos, bool) {
	op, opOffset, ok := p.lexer.ReadOp()
	if !ok {
		return 0, 0, false
	}
	return opToToken(op), p.offsetToPos(opOffset), true
}

func (p *Parser) withPrecedence(lhs ast.Expr, op token.Token, opPos token.Pos, rhs ast.Expr) *ast.BinaryExpr {
	binLhs, ok := lhs.(*ast.BinaryExpr)
	if !ok {
		return &ast.BinaryExpr{
			X:     lhs,
			Y:     rhs,
			Op:    op,
			OpPos: opPos,
		}
	}

	lExpr, mExpr, rExpr := binLhs.X, binLhs.Y, rhs
	lOp, rOp := binLhs.Op, op
	lOpPos, rOpPos := binLhs.OpPos, opPos

	if lOp.Precedence() >= rOp.Precedence() {
		lhs := &ast.BinaryExpr{
			X:     lExpr,
			Y:     mExpr,
			Op:    lOp,
			OpPos: lOpPos,
		}
		return &ast.BinaryExpr{
			X:     lhs,
			Y:     rExpr,
			Op:    rOp,
			OpPos: rOpPos,
		}
	}
	rhs = &ast.BinaryExpr{
		X:     mExpr,
		Y:     rExpr,
		Op:    rOp,
		OpPos: rOpPos,
	}
	return &ast.BinaryExpr{
		X:     lExpr,
		Y:     rhs,
		Op:    lOp,
		OpPos: lOpPos,
	}

}
