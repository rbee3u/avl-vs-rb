package trees

import (
	"errors"

	"github.com/rbee3u/avl-vs-rb/stats"
	"golang.org/x/exp/constraints"
)

var ErrUnimplemented = errors.New("unimplemented")

type Node[T constraints.Ordered] struct {
	parent *Node[T]
	left   *Node[T]
	right  *Node[T]
	extra  int8
	data   T
}

func NewNode[T constraints.Ordered](data T) *Node[T] {
	return &Node[T]{data: data}
}

func (n *Node[T]) Data() T {
	return n.data
}

func (n *Node[T]) next() *Node[T] {
	if n.right != nil {
		return minimum(n.right)
	}

	x := n

	for x == x.parent.right {
		x = x.parent
	}

	return x.parent
}

func (n *Node[T]) setExtra(extra int8) {
	stats.AddExtraCounter(1)

	n.extra = extra
}

func minimum[T constraints.Ordered](x *Node[T]) *Node[T] {
	for x.left != nil {
		x = x.left
	}

	return x
}

func maximum[T constraints.Ordered](x *Node[T]) *Node[T] {
	for x.right != nil {
		x = x.right
	}

	return x
}

func transplant[T constraints.Ordered](x *Node[T], y *Node[T]) {
	if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	if y != nil {
		y.parent = x.parent
	}
}

type Tree[T constraints.Ordered] interface {
	Size() int
	Empty() bool
	Begin() *Node[T]
	End() *Node[T]
	Clear()
	Find(data T) *Node[T]
	Insert(*Node[T])
	Delete(*Node[T])
}

type BaseTree[T constraints.Ordered] struct {
	sentinel Node[T]
	start    *Node[T]
	size     int
}

func NewBaseTree[T constraints.Ordered]() *BaseTree[T] {
	t := new(BaseTree[T])
	t.start = &t.sentinel

	return t
}

func (t *BaseTree[T]) Size() int {
	return t.size
}

func (t *BaseTree[T]) Empty() bool {
	return t.Size() == 0
}

func (t *BaseTree[T]) Begin() *Node[T] {
	return t.start
}

func (t *BaseTree[T]) End() *Node[T] {
	return &t.sentinel
}

func (t *BaseTree[T]) Clear() {
	t.End().left = nil
	t.start = t.End()
	t.size = 0
}

func (t *BaseTree[T]) Find(data T) *Node[T] {
	x := t.End()

	for y := x.left; y != nil; {
		stats.AddSearchCounter(1)

		switch {
		case data < y.data:
			y = y.left
		case y.data < data:
			y = y.right
		default:
			return y
		}
	}

	return x
}

func (t *BaseTree[T]) Insert(*Node[T]) {
	panic(ErrUnimplemented)
}

func (t *BaseTree[T]) Delete(*Node[T]) {
	panic(ErrUnimplemented)
}
