package ir

var (
	DisablePeephole = false
	nodeID          = 0
)

type Node interface {
	IsControl() bool
	base() *baseNode
}

type baseNode struct {
	ins  []Node
	outs []Node
	id   int
}

func initBaseNode[T Node](n T, ins ...Node) T {
	b := n.base()
	b.id = nodeID
	nodeID++
	b.ins = ins
	for _, in := range ins {
		if in != nil {
			addOut(in, n)
		}
	}
	return n
}

func In(n Node, i int) Node {
	return n.base().ins[i]
}

func NumIns(n Node) int {
	return len(n.base().ins)
}

func NumOuts(n Node) int {
	return len(n.base().outs)
}

func Unused(n Node) bool {
	return NumOuts(n) == 0
}

func Ins(n Node) []Node {
	return n.base().ins
}

func Outs(n Node) []Node {
	return n.base().outs
}

func addOut(n, out Node) {
	n.base().outs = append(n.base().outs, out)
}

func (b *baseNode) base() *baseNode {
	return b
}
