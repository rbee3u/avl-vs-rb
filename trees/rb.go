package trees

import (
	"github.com/rbee3u/avl-vs-rb/stats"
	"golang.org/x/exp/constraints"
)

const (
	rbRed   = 0
	rbBlack = 1
)

type RBTree[T constraints.Ordered] struct {
	*BaseTree[T]
}

func NewRBTree[T constraints.Ordered]() *RBTree[T] {
	return &RBTree[T]{BaseTree: NewBaseTree[T]()}
}

func (t *RBTree[T]) Insert(z *Node[T]) {
	z.parent = nil
	z.left, z.right, z.extra = nil, nil, rbRed
	x, childIsLeft := t.End(), true

	for y := x.left; y != nil; {
		stats.AddSearchCounter(1)

		x, childIsLeft = y, z.data < y.data

		if childIsLeft {
			y = y.left
		} else {
			y = y.right
		}
	}

	z.parent = x

	if childIsLeft {
		x.left = z
	} else {
		x.right = z
	}

	if t.start.left != nil {
		t.start = t.start.left
	}

	t.balanceAfterInsert(x, z)
	t.size++
}

func (t *RBTree[T]) balanceAfterInsert(x *Node[T], z *Node[T]) {
	for ; x != t.End() && x.extra == rbRed; x = z.parent {
		stats.AddFixupCounter(2)

		if x == x.parent.left {
			y := x.parent.right
			if rbIsRed(y) {
				z = z.parent
				z.setExtra(rbBlack)
				z = z.parent
				z.setExtra(rbRed)
				y.setExtra(rbBlack)
			} else {
				if z == x.right {
					z = x
					rbRotateLeft(z)
				}

				z = z.parent
				z.setExtra(rbBlack)
				z = z.parent
				z.setExtra(rbRed)
				rbRotateRight(z)
			}
		} else {
			y := x.parent.left
			if rbIsRed(y) {
				z = z.parent
				z.setExtra(rbBlack)
				z = z.parent
				z.setExtra(rbRed)
				y.setExtra(rbBlack)
			} else {
				if z == x.left {
					z = x
					rbRotateRight(z)
				}
				z = z.parent
				z.setExtra(rbBlack)
				z = z.parent
				z.setExtra(rbRed)
				rbRotateLeft(z)
			}
		}
	}

	t.End().left.extra = rbBlack
}

func (t *RBTree[T]) Delete(z *Node[T]) {
	if t.start == z {
		t.start = z.next()
	}

	x, deletedColor := z.parent, z.extra

	var n *Node[T]

	switch {
	case z.left == nil:
		n = z.right
		transplant(z, n)
	case z.right == nil:
		n = z.left
		transplant(z, n)
	default:
		y := minimum(z.right)
		x, deletedColor = y, y.extra
		n = y.right

		if y.parent != z {
			x = y.parent
			transplant(y, n)
			y.right = z.right
			y.right.parent = y
		}

		transplant(z, y)
		y.left = z.left
		y.left.parent = y
		y.extra = z.extra
	}

	if deletedColor == rbBlack {
		t.balanceAfterDelete(x, n)
	}
	t.size--
}

func (t *RBTree[T]) balanceAfterDelete(x *Node[T], n *Node[T]) {
	for ; x != t.End() && rbIsBlack(n); x = n.parent {
		stats.AddFixupCounter(1)

		if n == x.left {
			z := x.right
			if rbIsRed(z) {
				z.setExtra(rbBlack)
				x.setExtra(rbRed)
				rbRotateLeft(x)
				z = x.right
			}

			if rbIsBlack(z.left) && rbIsBlack(z.right) {
				z.setExtra(rbRed)

				n = x
			} else {
				if rbIsBlack(z.right) {
					z.left.setExtra(rbBlack)
					z.setExtra(rbRed)
					rbRotateRight(z)
					z = x.right
				}
				z.setExtra(x.extra)
				x.setExtra(rbBlack)
				z.right.setExtra(rbBlack)
				rbRotateLeft(x)
				n = t.End().left
			}
		} else {
			z := x.left
			if rbIsRed(z) {
				z.setExtra(rbBlack)
				x.setExtra(rbRed)
				rbRotateRight(x)
				z = x.left
			}
			if rbIsBlack(z.right) && rbIsBlack(z.left) {
				z.setExtra(rbRed)
				n = x
			} else {
				if rbIsBlack(z.left) {
					z.right.setExtra(rbBlack)
					z.setExtra(rbRed)
					rbRotateLeft(z)
					z = x.left
				}
				z.setExtra(x.extra)
				x.setExtra(rbBlack)
				z.left.setExtra(rbBlack)
				rbRotateRight(x)
				n = t.End().left
			}
		}
	}

	if rbIsRed(n) {
		n.setExtra(rbBlack)
	}
}

func rbRotateLeft[T constraints.Ordered](x *Node[T]) {
	stats.AddRotateCounter(1)

	y := x.right
	x.right = y.left

	if x.right != nil {
		x.right.parent = x
	}

	y.parent = x.parent

	if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	y.left = x
	x.parent = y
}

func rbRotateRight[T constraints.Ordered](x *Node[T]) {
	stats.AddRotateCounter(1)

	y := x.left
	x.left = y.right

	if x.left != nil {
		x.left.parent = x
	}

	y.parent = x.parent

	if x == x.parent.right {
		x.parent.right = y
	} else {
		x.parent.left = y
	}

	y.right = x
	x.parent = y
}

func rbIsRed[T constraints.Ordered](x *Node[T]) bool {
	return x != nil && x.extra == rbRed
}

func rbIsBlack[T constraints.Ordered](x *Node[T]) bool {
	return x == nil || x.extra == rbBlack
}
