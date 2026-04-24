package ir

import (
	"errors"
	"fmt"
	"go/ast"
	"slices"
	"strconv"
	"strings"

	"github.com/JackDalberg/SeaOfNodes/ir/types"
)

var (
	DisablePeephole = false
	nodeID          = 0
)

func computeError(n ast.Node, msg string) *ASTError {
	err := errors.New("Compute error: " + msg)
	return &ASTError{
		err: err,
		Pos: n.Pos(),
	}
}

type Node interface {
	IsControl() bool

	compute() (types.Type, error)
	label() string
	GraphicLabel() string
	toStringInternal(sb *strings.Builder)

	base() *baseNode
}

type baseNode struct {
	ins  []Node
	outs []Node
	id   int
	typ  types.Type
}

func (b *baseNode) base() *baseNode {
	return b
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

func UniqueName(n Node) string {
	id := strconv.Itoa(n.base().id)
	if _, ok := n.(*ConstantNode); ok {
		return "Con_" + id
	}
	return n.label() + id
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

func Type(n Node) types.Type {
	return n.base().typ
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

func removeOut(n, out Node) {
	n.base().outs = slices.DeleteFunc(n.base().outs, func(node Node) bool {
		return node == out
	})
}

func setIn(n Node, i int, in Node) error {
	old := In(n, i)
	if old == in {
		return nil
	}

	if old != nil {
		removeOut(old, n)
		if Unused(old) {
			err := kill(old)
			if err != nil {
				return err
			}
		}
	}

	n.base().ins[i] = in
	return nil
}

func kill(n Node) error {
	if !Unused(n) {
		return errors.New("Cannot kill a node that is in use")
	}

	for i := range n.base().ins {
		err := setIn(n, i, nil)
		if err != nil {
			return err
		}
	}

	n.base().ins = []Node{}
	n.base().typ = nil

	if !dead(n) {
		return fmt.Errorf("Node not dead after killing it: %v", n)
	}
	return nil
}

func dead(n Node) bool {
	return Unused(n) &&
		len(n.base().ins) == 0 &&
		n.base().typ == nil
}

func peephole(n Node) (Node, error) {
	typ, err := n.compute()
	if err != nil {
		return nil, err
	}
	n.base().typ = typ

	if DisablePeephole {
		return n, nil
	}

	if _, ok := n.(*ConstantNode); !ok && Type(n).Constant() {
		err := kill(n)
		if err != nil {
			return nil, err
		}
		return peephole(NewConstantNode(typ))
	}
	return n, nil
}

func ToString(n Node) string {
	sb := &strings.Builder{}
	toString(n, sb)
	return sb.String()
}

func toString(n Node, sb *strings.Builder) {
	if dead(n) {
		sb.WriteString(UniqueName(n))
		sb.WriteString(":DEAD")
		return
	}
	n.toStringInternal(sb)
}
