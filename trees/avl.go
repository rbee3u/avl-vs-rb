package trees

import (
	"github.com/rbee3u/avl-vs-rb/stats"
	"golang.org/x/exp/constraints"
)

const (
	avlLeftHeavy  = -1
	avlBalanced   = 0
	avlRightHeavy = +1
)

type AVLTree[T constraints.Ordered] struct {
	*BaseTree[T]
}

func NewAVLTree[T constraints.Ordered]() *AVLTree[T] {
	return &AVLTree[T]{BaseTree: NewBaseTree[T]()}
}

func (t *AVLTree[T]) Insert(z *Node[T]) {
	z.extra = avlBalanced
	z.parent, z.left, z.right = nil, nil, nil
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

	t.balanceAfterInsert(x, childIsLeft)
	t.size++
}

func (t *AVLTree[T]) balanceAfterInsert(x *Node[T], childIsLeft bool) {
	for ; x != t.End(); x = x.parent {
		stats.AddFixupCounter(1)

		if !childIsLeft {
			switch x.extra {
			case avlLeftHeavy:
				x.setExtra(avlBalanced)

				return
			case avlRightHeavy:
				if x.right.extra == avlLeftHeavy {
					avlRotateRightLeft(x)
				} else {
					avlRotateLeft(x)
				}

				return
			default:
				x.setExtra(avlRightHeavy)
			}
		} else {
			switch x.extra {
			case avlRightHeavy:
				x.setExtra(avlBalanced)

				return
			case avlLeftHeavy:
				if x.left.extra == avlRightHeavy {
					avlRotateLeftRight(x)
				} else {
					avlRotateRight(x)
				}

				return
			default:
				x.setExtra(avlLeftHeavy)
			}
		}

		childIsLeft = x == x.parent.left
	}
}

func (t *AVLTree[T]) Delete(z *Node[T]) {
	if t.start == z {
		t.start = z.next()
	}

	x, childIsLeft := z.parent, z == z.parent.left

	switch {
	case z.left == nil:
		transplant(z, z.right)
	case z.right == nil:
		transplant(z, z.left)
	default:
		if z.extra == avlRightHeavy {
			y := minimum(z.right)
			x, childIsLeft = y, y == y.parent.left

			if y.parent != z {
				x = y.parent
				transplant(y, y.right)
				y.right = z.right
				y.right.parent = y
			}

			transplant(z, y)
			y.left = z.left
			y.left.parent = y
			y.extra = z.extra
		} else {
			y := maximum(z.left)
			x, childIsLeft = y, y == y.parent.left

			if y.parent != z {
				x = y.parent
				transplant(y, y.left)
				y.left = z.left
				y.left.parent = y
			}

			transplant(z, y)
			y.right = z.right
			y.right.parent = y
			y.extra = z.extra
		}
	}

	t.balanceAfterDelete(x, childIsLeft)
	t.size--
}

func (t *AVLTree[T]) balanceAfterDelete(x *Node[T], childIsLeft bool) {
	for ; x != t.End(); x = x.parent {
		stats.AddFixupCounter(1)

		if childIsLeft {
			switch x.extra {
			case avlBalanced:
				x.setExtra(avlRightHeavy)

				return
			case avlRightHeavy:
				b := x.right.extra
				if b == avlLeftHeavy {
					avlRotateRightLeft(x)
				} else {
					avlRotateLeft(x)
				}

				if b == avlBalanced {
					return
				}

				x = x.parent
			default:
				x.setExtra(avlBalanced)
			}
		} else {
			switch x.extra {
			case avlBalanced:
				x.setExtra(avlLeftHeavy)

				return
			case avlLeftHeavy:
				b := x.left.extra
				if b == avlRightHeavy {
					avlRotateLeftRight(x)
				} else {
					avlRotateRight(x)
				}
				if b == avlBalanced {
					return
				}
				x = x.parent
			default:
				x.setExtra(avlBalanced)
			}
		}

		childIsLeft = x == x.parent.left
	}
}

func avlRotateLeft[T constraints.Ordered](x *Node[T]) {
	stats.AddRotateCounter(1)

	z := x.right
	x.right = z.left

	if z.left != nil {
		z.left.parent = x
	}

	z.parent = x.parent

	if x == x.parent.left {
		x.parent.left = z
	} else {
		x.parent.right = z
	}

	z.left = x
	x.parent = z

	if z.extra == avlBalanced {
		x.setExtra(avlRightHeavy)
		z.setExtra(avlLeftHeavy)
	} else {
		x.setExtra(avlBalanced)
		z.setExtra(avlBalanced)
	}
}

func avlRotateRight[T constraints.Ordered](x *Node[T]) {
	stats.AddRotateCounter(1)

	z := x.left
	x.left = z.right

	if z.right != nil {
		z.right.parent = x
	}

	z.parent = x.parent

	if x == x.parent.right {
		x.parent.right = z
	} else {
		x.parent.left = z
	}

	z.right = x
	x.parent = z

	if z.extra == avlBalanced {
		x.setExtra(avlLeftHeavy)
		z.setExtra(avlRightHeavy)
	} else {
		x.setExtra(avlBalanced)
		z.setExtra(avlBalanced)
	}
}

func avlRotateRightLeft[T constraints.Ordered](x *Node[T]) {
	stats.AddRotateCounter(2)

	z := x.right
	y := z.left
	z.left = y.right

	if y.right != nil {
		y.right.parent = z
	}

	y.right = z
	z.parent = y
	x.right = y.left

	if y.left != nil {
		y.left.parent = x
	}

	y.parent = x.parent

	if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	y.left = x
	x.parent = y

	switch y.extra {
	case avlRightHeavy:
		x.setExtra(avlLeftHeavy)
		y.setExtra(avlBalanced)
		z.setExtra(avlBalanced)
	case avlLeftHeavy:
		x.setExtra(avlBalanced)
		y.setExtra(avlBalanced)
		z.setExtra(avlRightHeavy)
	default:
		x.setExtra(avlBalanced)
		z.setExtra(avlBalanced)
	}
}

func avlRotateLeftRight[T constraints.Ordered](x *Node[T]) {
	stats.AddRotateCounter(2)

	z := x.left
	y := z.right
	z.right = y.left

	if y.left != nil {
		y.left.parent = z
	}

	y.left = z
	z.parent = y
	x.left = y.right

	if y.right != nil {
		y.right.parent = x
	}

	y.parent = x.parent

	if x == x.parent.right {
		x.parent.right = y
	} else {
		x.parent.left = y
	}

	y.right = x
	x.parent = y

	switch y.extra {
	case avlLeftHeavy:
		x.setExtra(avlRightHeavy)
		y.setExtra(avlBalanced)
		z.setExtra(avlBalanced)
	case avlRightHeavy:
		x.setExtra(avlBalanced)
		y.setExtra(avlBalanced)
		z.setExtra(avlLeftHeavy)
	default:
		x.setExtra(avlBalanced)
		z.setExtra(avlBalanced)
	}
}
