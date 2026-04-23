package ir

// Only need 1 start node
var StartNode = newStartNode()

type startNode struct {
	baseNode
}

func newStartNode() *startNode {
	return initBaseNode(&startNode{})
}

func (s *startNode) IsControl() bool {
	return true
}
