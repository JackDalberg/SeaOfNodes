package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"

	"github.com/JackDalberg/SeaOfNodes/ir"
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

func (p *Parser) Parse() (ast.Node, error) {
	n, err := p.parseBlock(p.offsetToPos(0), false)
	if err != nil {
		return nil, err
	}
	if b, offset, ok := p.lexer.ReadByte(); ok {
		return nil, syntaxError(offset, "unexpected %c", b)
	}
	return n, nil
}

func (p *Parser) parseBlock(pos token.Pos, endInCurly bool) (*ast.BlockStmt, error) {
	block := &ast.BlockStmt{Lbrace: pos}
	for !p.blockEnd(block, endInCurly) {
		n, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if n != nil {
			block.List = append(block.List, n)
		}
	}
	return block, nil
}

func (p *Parser) blockEnd(block *ast.BlockStmt, endInCurly bool) bool {
	if !endInCurly {
		return p.lexer.isEOF()
	}

	offset, ok := p.lexer.Read('}')
	if !ok {
		return false
	}
	block.Rbrace = p.offsetToPos(offset)
	return true
}

func (p *Parser) parseStatement() (ast.Stmt, error) {
	t, offset, isID, err := p.lexer.ReadToken()
	if err != nil {
		return nil, err
	}

	pos := p.offsetToPos(offset)
	var n ast.Stmt
	switch t {
	case "return":
		n, err = p.parseReturn(pos)
		if err != nil {
			return nil, err
		}

	case "int":
		n, err = p.parseDecl(pos)
		if err != nil {
			return nil, err
		}

	case "#":
		return p.parseInstruction()

	case "{":
		n, err = p.parseBlock(pos, true)
		if err != nil {
			return nil, err
		}

	case ";":
		return nil, nil

	default:
		if isID {
			return p.parseExprStatement(t, pos)
		}
		return nil, syntaxError(offset, "expected a statement, got=%q", t)
	}

	return n, nil
}

func (p *Parser) parseInstruction() (ast.Stmt, error) {
	inst, offset, _, err := p.lexer.ReadToken()
	if err != nil {
		return nil, err
	}

	switch inst {
	case "showGraph":
		return ir.ShowGraphInst, nil

	case "disablePeephole":
		return ir.DisablePeepholeInst, nil
	}
	return nil, syntaxError(offset, "unknown compiler instruction")
}

func (p *Parser) parseReturn(pos token.Pos) (*ast.ReturnStmt, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	err = p.parseSemicolon()
	if err != nil {
		return nil, err
	}
	return &ast.ReturnStmt{Return: pos, Results: []ast.Expr{expr}}, nil
}

func (p *Parser) parseExprStatement(name string, namePos token.Pos) (ast.Stmt, error) {
	id := &ast.Ident{NamePos: namePos, Name: name}
	offset, ok := p.lexer.Read('=')
	if !ok {
		return nil, syntaxError(offset, "expected assignment '='")
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	assign := &ast.AssignStmt{
		Lhs:    []ast.Expr{id},
		Tok:    token.EQL,
		TokPos: p.offsetToPos(offset),
		Rhs:    []ast.Expr{expr},
	}
	return assign, nil
}

func (p *Parser) parseSemicolon() error {
	offset, ok := p.lexer.Read(';')
	if !ok {
		return syntaxError(offset, "expected ';' after expression")
	}
	return nil
}

func (p *Parser) parseID() (*ast.Ident, error) {
	name, nameOffset, ok := p.lexer.ReadID()
	if !ok {
		return nil, syntaxError(nameOffset, "expected identifier")
	}
	return &ast.Ident{NamePos: p.offsetToPos(nameOffset), Name: name}, nil
}

func (p *Parser) parseDecl(pos token.Pos) (*ast.DeclStmt, error) {
	name, err := p.parseID()
	if err != nil {
		return nil, err
	}

	opOffset, ok := p.lexer.Read('=')
	if !ok {
		return nil, syntaxError(opOffset, "expected '='")
	}

	value, err := p.parseBinary()
	if err != nil {
		return nil, err
	}

	err = p.parseSemicolon()
	if err != nil {
		return nil, err
	}

	decl := &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names:  []*ast.Ident{name},
					Type:   &ast.Ident{Name: "int", NamePos: pos},
					Values: []ast.Expr{value},
				},
			},
		},
	}
	return decl, nil
}

func (p *Parser) parseExpr() (ast.Expr, error) {
	return p.parseBinary()
}

func (p *Parser) parsePrimary() (ast.Expr, error) {
	num, offset, err := p.lexer.ReadNumber()
	if err != nil {
		if errors.Is(err, ErrNAN) {
			return p.parseID()
		}
		return nil, syntaxError(offset, err.Error())
	}
	blit := &ast.BasicLit{
		ValuePos: p.offsetToPos(offset),
		Kind:     token.INT,
		Value:    num,
	}
	return blit, nil
}

func (p *Parser) parseUnary() (ast.Expr, error) {
	if lOffset, ok := p.lexer.Read('('); ok {
		expr, err := p.parseBinary()
		if err != nil {
			return nil, err
		}
		rOffset, ok := p.lexer.Read(')')
		if !ok {
			return nil, syntaxError(rOffset, "expected ')'")
		}
		paren := &ast.ParenExpr{
			Lparen: p.offsetToPos(lOffset),
			X:      expr,
			Rparen: p.offsetToPos(rOffset),
		}
		return paren, nil
	}

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

func (p *Parser) offsetToPos(offset int) token.Pos {
	return p.file.Pos(offset)
}

func (p *Parser) PosToOffset(pos token.Pos) int {
	return p.file.Offset(pos)
}

func (p *Parser) string(n ast.Node) string {
	sb := &strings.Builder{}
	printer.Fprint(sb, p.fileset, n)
	return sb.String()
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
	case '=':
		return token.EQL
	default:
		return token.ILLEGAL
	}
}
