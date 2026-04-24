package ir

import (
	"seaofnodes/ir/types"
	"strings"
)

// Only need 1 start node
var StartNode = newStartNode()

type startNode struct {
	baseNode
}

func (s *startNode) compute() (types.Type, error) {
	return types.BottomType, nil
}

func (s *startNode) label() string {
	return "Start"
}

func (s *startNode) GraphicLabel() string {
	return s.label()
}

func (s *startNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString(s.label())
}

func newStartNode() *startNode {
	return initBaseNode(&startNode{})
}

func (s *startNode) IsControl() bool {
	return true
}
