package main

import (
	"errors"
	goParser "go/parser"
	"seaofnodes/ir"
	"seaofnodes/parser"
	"strings"
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

func Simple(source string) (*ir.ReturnNode, error) {
	p := parser.NewParser(source)
	n, err := p.Parse()
	if err != nil {
		if s, ok := errors.AsType[*parser.SyntaxError](err); ok {
			return nil, &SourceError{err: s, source: source, offset: s.Offset}
		}
		return nil, err
	}
	generator := ir.NewGenerator()
	ret, err := generator.Generate(n)
	if err != nil {
		if a, ok := errors.AsType[*ir.ASTError](err); ok {
			return nil, &SourceError{err: a, source: source, offset: p.PosToOffset(a.Pos)}
		}
		return nil, err
	}
	return ret, nil
}

func GoSimple(source string) (*ir.ReturnNode, error) {
	n, err := goParser.ParseExpr(source)
	if err != nil {
		return nil, err
	}
	generator := ir.NewGenerator()
	return generator.Generate(n)
}
