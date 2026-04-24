package simple

import (
	"errors"
	goParser "go/parser"
	"strings"

	"github.com/JackDalberg/SeaOfNodes/ir"
	"github.com/JackDalberg/SeaOfNodes/parser"
)

type SourceError struct {
	err    error
	source string
	offset int
}

func (s *SourceError) Error() string {
	msg := "\n" + s.source + "\n"
	msg += strings.Repeat(" ", s.offset) + "^\n"
	msg += s.err.Error()
	return msg
}

func Simple(source string) (*ir.ReturnNode, *ir.Generator, error) {
	p := parser.NewParser(source)
	n, err := p.Parse()
	if err != nil {
		if s, ok := errors.AsType[*parser.SyntaxError](err); ok {
			return nil, nil, &SourceError{err: s, source: source, offset: s.Offset}
		}
		return nil, nil, err
	}
	generator := ir.NewGenerator()
	ret, err := generator.Generate(n)
	if err != nil {
		if a, ok := errors.AsType[*ir.ASTError](err); ok {
			return nil, nil, &SourceError{err: a, source: source, offset: p.PosToOffset(a.Pos)}
		}
		return nil, nil, err
	}
	return ret, generator, nil
}

func GoSimple(source string) (*ir.ReturnNode, *ir.Generator, error) {
	n, err := goParser.ParseExpr(source)
	if err != nil {
		return nil, nil, err
	}
	generator := ir.NewGenerator()
	retNode, err := generator.Generate(n)
	if err != nil {
		return nil, nil, err
	}
	return retNode, generator, nil
}
