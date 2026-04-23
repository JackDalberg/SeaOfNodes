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
	return p.parsePrimary()
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
