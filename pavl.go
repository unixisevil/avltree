package avltree

import (
	"unsafe"
)

type pnode struct {
	links   [ChildNum]*pnode //child node
	parent  *pnode           //parent node
	data    Item             //data item
	balance int8             //balance factor
}

type PAvlTree struct {
	root       *pnode      //root of  tree
	cmpFunc    Compare     //compare function
	extraParam interface{} //extra param for cmpFunc
	count      int         // number of item in tree
}

func NewPAvl(cmp Compare, extra interface{}) *PAvlTree {
	if cmp == nil {
		return nil
	}
	return &PAvlTree{
		cmpFunc:    cmp,
		extraParam: extra,
	}
}

func (t *PAvlTree) Count() int {
	if t == nil {
		return 0
	}
	return t.count
}

//search target in tree
//if find it return item
//else return nil
func (t *PAvlTree) Find(target Item) Item {
	if t == nil || target == nil {
		return nil
	}
	for w := t.root; w != nil; {
		ret := t.cmpFunc(target, w.data, t.extraParam)
		if ret < 0 {
			w = w.links[Left]
		} else if ret > 0 {
			w = w.links[Right]
		} else {
			return w.data
		}
	}
	return nil
}

func (t *PAvlTree) insert(item Item) (*Item, bool) {
	if t == nil || item == nil {
		return nil, false
	}
	var (
		y   *pnode //待更新平衡因子的最顶层节点
		w   *pnode //current walk node
		p   *pnode //w's  parent
		n   *pnode //new node
		r   *pnode //new root node of rebalanced subtree
		dir byte   //下降方向
	)
	y = t.root
	for p, w = nil, t.root; w != nil; p, w = w, w.links[dir] {
		cmp := t.cmpFunc(item, w.data, t.extraParam)
		if cmp == 0 {
			return &w.data, false
		}
		if cmp > 0 {
			dir = Right
		} else {
			dir = Left
		}
		if w.balance != 0 {
			y = w
		}
	}
	n = &pnode{data: item, parent: p}
	t.count++
	if p != nil {
		p.links[dir] = n
	} else {
		t.root = n
	}
	if t.root == n {
		return &n.data, true
	}
	for w = n; w != y; w = p {
		p = w.parent
		if p.links[Left] != w {
			p.balance++
		} else {
			p.balance--
		}
	}
	if y.balance == -2 {
		//fmt.Printf("y.balance == -2\n")
		x := y.links[Left]
		if x.balance == -1 {
			r = x
			y.links[Left] = x.links[Right]
			x.links[Right] = y
			x.balance = 0
			y.balance = 0
			x.parent = y.parent
			y.parent = x
			if y.links[Left] != nil {
				y.links[Left].parent = y
			}
		} else { //x.balance == 1
			r = x.links[Right]
			x.links[Right] = r.links[Left]
			r.links[Left] = x
			y.links[Left] = r.links[Right]
			r.links[Right] = y
			if r.balance == -1 {
				x.balance = 0
				y.balance = 1
			} else if r.balance == 0 {
				x.balance = 0
				y.balance = 0
			} else {
				x.balance = -1
				y.balance = 0
			}
			r.balance = 0
			r.parent = y.parent
			x.parent = r
			y.parent = r
			if x.links[Right] != nil {
				x.links[Right].parent = x
			}
			if y.links[Left] != nil {
				y.links[Left].parent = y
			}
		}
	} else if y.balance == 2 {
		//fmt.Printf("y.balance == 2\n")
		x := y.links[Right]
		if x.balance == 1 {
			r = x
			y.links[Right] = x.links[Left]
			x.links[Left] = y
			x.balance = 0
			y.balance = 0
			x.parent = y.parent
			y.parent = x
			if y.links[Right] != nil {
				y.links[Right].parent = y
			}
		} else { //x->avl_balance == -1
			r = x.links[Left]
			x.links[Left] = r.links[Right]
			r.links[Right] = x
			y.links[Right] = r.links[Left]
			r.links[Left] = y
			if r.balance == 1 {
				x.balance = 0
				y.balance = -1
			} else if r.balance == 0 {
				x.balance = 0
				y.balance = 0
			} else {
				x.balance = 1
				y.balance = 0
			}
			r.balance = 0
			r.parent = y.parent
			x.parent = r
			y.parent = r
			if x.links[Left] != nil {
				x.links[Left].parent = x
			}
			if y.links[Right] != nil {
				y.links[Right].parent = y
			}
		}
	} else {
		//fmt.Printf("no balance return\n")
		return &n.data, true
	}
	if r.parent != nil {
		if r.parent.links[Left] != y {
			dir = Right
		} else {
			dir = Left
		}
		r.parent.links[dir] = r
	} else {
		t.root = r
	}
	return &n.data, true
}

//insert item in tree
//return true if item was successfully inserted
//return false if item already in tree
func (t *PAvlTree) Insert(item Item) bool {
	_, succ := t.insert(item)
	return succ
}

//replace item in tree with same key item
//return old item
func (t *PAvlTree) Replace(item Item) Item {
	addr, succ := t.insert(item)
	if addr == nil || succ {
		return nil
	}
	r := *addr
	*addr = item
	return r
}

//delete item in tree
//return item if find it
//else  return nil
func (t *PAvlTree) Delete(item Item) Item {
	if t == nil || item == nil {
		return nil
	}
	if t.root == nil {
		return nil
	}
	var dir int
	w := t.root //walk node
	for {
		cmp := t.cmpFunc(item, w.data, t.extraParam)
		if cmp == 0 {
			break
		}
		if cmp > 0 {
			dir = Right
		} else {
			dir = Left
		}
		w = w.links[dir]
		if w == nil {
			return nil
		}
	}
	ret := w.data
	p := w.parent
	if p == nil {
		p = (*pnode)(unsafe.Pointer(&t.root))
		dir = Left
	}
	if w.links[Right] == nil { //case 1, w has no right child
		p.links[dir] = w.links[Left]
		if p.links[dir] != nil {
			p.links[dir].parent = w.parent
		}
	} else { //case 2, w's right child has no left child
		r := w.links[Right]
		if r.links[Left] == nil {
			r.links[Left] = w.links[Left]
			p.links[dir] = r
			r.parent = w.parent
			if r.links[Left] != nil {
				r.links[Left].parent = r
			}
			r.balance = w.balance
			p = r
			dir = Right
		} else { //case 3, w's right child has left child
			s := r.links[Left]
			for s.links[Left] != nil {
				s = s.links[Left]
			}
			r = s.parent
			r.links[Left] = s.links[Right]
			s.links[Left] = w.links[Left]
			s.links[Right] = w.links[Right]
			p.links[dir] = s
			if s.links[Left] != nil {
				s.links[Left].parent = s
			}
			s.links[Right].parent = s
			s.parent = w.parent
			if r.links[Left] != nil {
				r.links[Left].parent = r
			}
			s.balance = w.balance
			p = r
			dir = Left
		}
	}
	w = nil
	for p != (*pnode)(unsafe.Pointer(&t.root)) {
		y := p
		if y.parent != nil {
			p = y.parent
		} else {
			p = (*pnode)(unsafe.Pointer(&t.root))
		}
		if dir == Left {
			if p.links[Left] != y {
				dir = Right
			} else {
				dir = Left
			}
			y.balance++
			if y.balance == 1 {
				break
			} else if y.balance == 2 {
				x := y.links[Right]
				if x.balance == -1 {
					r := x.links[Left]
					x.links[Left] = r.links[Right]
					r.links[Right] = x
					y.links[Right] = r.links[Left]
					r.links[Left] = y
					if r.balance == 1 {
						x.balance = 0
						y.balance = -1
					} else if r.balance == 0 {
						x.balance = 0
						y.balance = 0
					} else { /* r.balance == -1 */
						x.balance = 1
						y.balance = 0
					}
					r.balance = 0
					r.parent = y.parent
					x.parent = r
					y.parent = r
					if x.links[Left] != nil {
						x.links[Left].parent = x
					}
					if y.links[Right] != nil {
						y.links[Right].parent = y
					}
					p.links[dir] = r
				} else { /*  x.balance == 0  ||  x.balance == 1 */
					y.links[Right] = x.links[Left]
					x.links[Left] = y
					x.parent = y.parent
					y.parent = x
					if y.links[Right] != nil {
						y.links[Right].parent = y
					}
					p.links[dir] = x
					if x.balance == 0 {
						x.balance = -1
						y.balance = 1
						break
					} else {
						x.balance = 0
						y.balance = 0
						y = x
					}
				}
			}

		} else { // dir == Right
			if p.links[Left] != y {
				dir = Right
			} else {
				dir = Left
			}
			y.balance--
			if y.balance == -1 {
				break
			} else if y.balance == -2 {
				x := y.links[Left]
				if x.balance == 1 {
					r := x.links[Right]
					x.links[Right] = r.links[Left]
					r.links[Left] = x
					y.links[Left] = r.links[Right]
					r.links[Right] = y
					if r.balance == -1 {
						x.balance = 0
						y.balance = 1
					} else if r.balance == 0 {
						x.balance = 0
						y.balance = 0
					} else {
						x.balance = -1
						y.balance = 0
					}
					r.balance = 0
					r.parent = y.parent
					x.parent = r
					y.parent = r
					if x.links[Right] != nil {
						x.links[Right].parent = x
					}
					if y.links[Left] != nil {
						y.links[Left].parent = y
					}
					p.links[dir] = r
				} else {
					y.links[Left] = x.links[Right]
					x.links[Right] = y
					x.parent = y.parent
					y.parent = x
					if y.links[Left] != nil {
						y.links[Left].parent = y
					}
					p.links[dir] = x
					if x.balance == 0 {
						x.balance = 1
						y.balance = -1
						break
					} else {
						x.balance = 0
						y.balance = 0
						y = x
					}
				}
			}
		}
	}
	t.count--
	return ret
}

func (t *PAvlTree) Copy() *PAvlTree {
	if t == nil {
		return nil
	}
	n := NewPAvl(t.cmpFunc, t.extraParam)
	if n == nil {
		return nil
	}
	n.count = t.count
	if n.count == 0 {
		return n
	}
	var (
		x *pnode
		y *pnode
	)
	x = (*pnode)(unsafe.Pointer(&t.root))
	y = (*pnode)(unsafe.Pointer(&n.root))
	for {
		for x.links[Left] != nil {
			y.links[Left] = &pnode{}
			y.links[Left].parent = y
			x = x.links[Left]
			y = y.links[Left]
		}
		y.links[Left] = nil
		for {
			y.data = x.data
			y.balance = x.balance
			if x.links[Right] != nil {
				y.links[Right] = &pnode{}
				y.links[Right].parent = y
				x = x.links[Right]
				y = y.links[Right]
				break
			} else {
				y.links[Right] = nil
			}
			for {
				w := x
				x = x.parent
				if x == nil {
					n.root.parent = nil
					return n
				}
				y = y.parent
				if w == x.links[Left] {
					break
				}
			}
		}
	}
}

func (t *PAvlTree) Iter() *PAvlIter {
	it := NewPIter()
	return it.HookWith(t)
}

type PAvlIter struct {
	tree *PAvlTree //the tree be iterated
	node *pnode    //current node in tree
}

func NewPIter() *PAvlIter {
	return &PAvlIter{}
}

func (it *PAvlIter) HookWith(tree *PAvlTree) *PAvlIter {
	if it == nil {
		return nil
	}
	it.tree = tree
	it.node = nil
	return it
}

func (it *PAvlIter) First() Item {
	if it == nil || it.tree == nil {
		return nil
	}
	w := it.tree.root
	if w == nil {
		return nil
	}
	for w.links[Left] != nil {
		w = w.links[Left]
	}
	it.node = w
	return w.data
}

func (it *PAvlIter) Last() Item {
	if it == nil || it.tree == nil {
		return nil
	}
	w := it.tree.root
	if w == nil {
		return nil
	}
	for w.links[Right] != nil {
		w = w.links[Right]
	}
	it.node = w
	return w.data
}

func (it *PAvlIter) Find(item Item) Item {
	if it == nil || it.tree == nil || item == nil {
		return nil
	}
	var (
		w *pnode //walk node
		n *pnode //child of w
	)
	for w = it.tree.root; w != nil; w = n {
		cmp := it.tree.cmpFunc(item, w.data, it.tree.extraParam)
		if cmp == 0 {
			it.node = w
			return w.data
		}
		if cmp < 0 {
			n = w.links[Left]
		} else {
			n = w.links[Right]
		}
	}
	it.node = nil
	return nil
}

func (it *PAvlIter) Next() Item {
	if it == nil || it.tree == nil {
		return nil
	}
	w := it.node
	if w == nil {
		return it.First()
	} else if w.links[Right] != nil {
		w = w.links[Right]
		for w.links[Left] != nil {
			w = w.links[Left]
		}
	} else {
		for p := w.parent; ; w, p = p, p.parent {
			if p == nil {
				it.node = p
				return nil
			}
			if w == p.links[Left] {
				it.node = p
				return p.data
			}
		}
	}
	it.node = w
	return w.data
}

func (it *PAvlIter) Prev() Item {
	if it == nil || it.tree == nil {
		return nil
	}
	w := it.node
	if w == nil {
		return it.Last()
	} else if w.links[Left] != nil {
		w = w.links[Left]
		for w.links[Right] != nil {
			w = w.links[Right]
		}
	} else {
		for p := w.parent; ; w, p = p, p.parent {
			if p == nil {
				it.node = p
				return nil
			}
			if w == p.links[Right] {
				it.node = p
				return p.data
			}
		}
	}
	it.node = w
	return w.data
}

func (it *PAvlIter) Current() Item {
	if it == nil || it.node == nil {
		return nil
	}
	return it.node.data
}

//don't change key part of item
func (it *PAvlIter) Replace(new Item) Item {
	if it == nil || it.node == nil || new == nil {
		return nil
	}
	old := it.node.data
	it.node.data = new
	return old
}

func (it *PAvlIter) CopyFrom(other *PAvlIter) Item {
	if it == nil || other == nil {
		return nil
	}
	it.tree = other.tree
	it.node = other.node
	if it.node == nil {
		return nil
	}
	return it.node.data
}

func (it *PAvlIter) Insert(item Item) (*Item, bool) {
	if it == nil || it.tree == nil || item == nil {
		return nil, false
	}
	addr, ok := it.tree.insert(item)
	it.node = (*pnode)(unsafe.Pointer(uintptr(unsafe.Pointer(addr)) - unsafe.Offsetof(it.node.data)))
	return addr, ok
}
