package bbst

import (
	"unsafe"
)

type prbnode struct {
	links  [ChildNum]*prbnode //child node
	parent *prbnode           //parent node
	data   Item               //data item
	color  byte               //node color
}

type PRbTree struct {
	root       *prbnode    //root of  tree
	cmpFunc    Compare     //compare function
	extraParam interface{} //extra param for cmpFunc
	count      int         // number of item in tree
}

func NewPRbTree(cmp Compare, extra interface{}) *PRbTree {
	if cmp == nil {
		return nil
	}
	return &PRbTree{
		cmpFunc:    cmp,
		extraParam: extra,
	}
}

func (t *PRbTree) Count() int {
	if t == nil {
		return 0
	}
	return t.count
}

//search target in tree
//if find it return item
//else return nil
func (t *PRbTree) Find(target Item) Item {
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

func (t *PRbTree) insert(item Item) (*Item, bool) {
	if t == nil || item == nil {
		return nil, false
	}
	var (
		w   *prbnode //walk node
		p   *prbnode //parent of w
		n   *prbnode //new node
		dir int      //direction of p
	)
	for w = t.root; w != nil; p, w = w, w.links[dir] {
		cmp := t.cmpFunc(item, w.data, t.extraParam)
		if cmp == 0 {
			return &w.data, false
		}
		if cmp > 0 {
			dir = Right
		} else {
			dir = Left
		}
	}
	n = &prbnode{parent: p, data: item, color: red}
	if p != nil {
		p.links[dir] = n
	} else {
		t.root = n
	}
	t.count++

	w = n
	for {
		p = w.parent
		if p == nil || p.color == black {
			break
		}
		g := p.parent
		if g == nil {
			break
		}
		if g.links[Left] == p {
			y := g.links[Right]
			if y != nil && y.color == red {
				p.color = black
				y.color = black
				g.color = red
				w = g
			} else {
				pg := g.parent
				if pg == nil {
					pg = (*prbnode)(unsafe.Pointer(&t.root))
				}
				if p.links[Right] == w {
					p.links[Right] = w.links[Left]
					w.links[Left] = p
					g.links[Left] = w
					p.parent = w
					if p.links[Right] != nil {
						p.links[Right].parent = p
					}
					p = w
				}
				g.color = red
				p.color = black
				g.links[Left] = p.links[Right]
				p.links[Right] = g

				d := Left
				if pg.links[Left] != g {
					d = Right
				}
				pg.links[d] = p
				p.parent = g.parent
				g.parent = p
				if g.links[Left] != nil {
					g.links[Left].parent = g
				}
				break
			}
		} else {
			y := g.links[Left]
			if y != nil && y.color == red {
				p.color = black
				y.color = black
				g.color = red
				w = g
			} else {
				pg := g.parent
				if pg == nil {
					pg = (*prbnode)(unsafe.Pointer(&t.root))
				}
				if p.links[Left] == w {
					p.links[Left] = w.links[Right]
					w.links[Right] = p
					g.links[Right] = w
					p.parent = w
					if p.links[Left] != nil {
						p.links[Left].parent = p
					}
					p = w
				}
				g.color = red
				p.color = black
				g.links[Right] = p.links[Left]
				p.links[Left] = g

				d := Left
				if pg.links[Left] != g {
					d = Right
				}
				pg.links[d] = p

				p.parent = g.parent
				g.parent = p
				if g.links[Right] != nil {
					g.links[Right].parent = g
				}
				break
			}
		}
	}
	t.root.color = black
	return &n.data, true
}

//insert item in tree
//return true if item was successfully inserted
//return false if item already in tree
func (t *PRbTree) Insert(item Item) bool {
	_, succ := t.insert(item)
	return succ
}

//replace item in tree with same key item
func (t *PRbTree) Replace(item Item) Item {
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
func (t *PRbTree) Delete(item Item) Item {
	if t == nil || item == nil {
		return nil
	}
	if t.root == nil {
		return nil
	}
	var (
		w   *prbnode //node to delete
		p   *prbnode //parent of w
		f   *prbnode //rebalancing node
		dir int      //direction of p or f
	)
	for w = t.root; ; {
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
	p = w.parent
	if p == nil {
		p = (*prbnode)(unsafe.Pointer(&t.root))
		dir = Left
	}
	if w.links[Right] == nil { //case 1, node to delete has no right child
		p.links[dir] = w.links[Left]
		if p.links[dir] != nil {
			p.links[dir].parent = w.parent
		}
		//rebalancing start at p
		f = p
	} else {
		r := w.links[Right]
		if r.links[Left] == nil { //case 2, node to delete w's right child has no left child
			r.links[Left] = w.links[Left]
			p.links[dir] = r
			r.parent = w.parent
			if r.links[Left] != nil {
				r.links[Left].parent = r
			}
			w.color, r.color = r.color, w.color
			f = r
			dir = Right
		} else { //case 3, node to delete w's right child has left child
			s := r.links[Left]
			for s.links[Left] != nil {
				s = s.links[Left]
			}
			r = s.parent
			//cut off s from r, replace w with s
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
			w.color, s.color = s.color, w.color
			f = r
			dir = Left
		}
	}
	if w.color == black {
		for {
			var tmp *prbnode
			x := f.links[dir]
			if x != nil && x.color == red {
				x.color = black
				break
			}
			if f == (*prbnode)(unsafe.Pointer(&t.root)) {
				break
			}
			g := f.parent
			if g == nil {
				g = (*prbnode)(unsafe.Pointer(&t.root))
			}
			if dir == Left {
				//node x's sibling
				s := f.links[Right]
				if s.color == red {
					s.color = black
					f.color = red
					f.links[Right] = s.links[Left]
					s.links[Left] = f

					d := Left
					if g.links[Left] != f {
						d = Right
					}
					g.links[d] = s

					s.parent = f.parent
					f.parent = s

					g = s
					s = f.links[Right]
					s.parent = f
				}
				if (s.links[Left] == nil || s.links[Left].color == black) &&
					(s.links[Right] == nil || s.links[Right].color == black) {
					s.color = red
				} else {
					if s.links[Right] == nil || s.links[Right].color == black {
						y := s.links[Left]
						y.color = black
						s.color = red
						s.links[Left] = y.links[Right]
						y.links[Right] = s
						if s.links[Left] != nil {
							s.links[Left].parent = s
						}
						f.links[Right] = y
						s = y
						s.links[Right].parent = s
					}
					s.color = f.color
					f.color = black
					s.links[Right].color = black

					f.links[Right] = s.links[Left]
					s.links[Left] = f

					d := Left
					if g.links[Left] != f {
						d = Right
					}
					g.links[d] = s

					s.parent = f.parent
					f.parent = s
					if f.links[Right] != nil {
						f.links[Right].parent = f
					}
					break
				}
			} else { // dir == Right
				//node x's sibling
				s := f.links[Left]
				if s.color == red {
					s.color = black
					f.color = red
					f.links[Left] = s.links[Right]
					s.links[Right] = f

					d := Left
					if g.links[Left] != f {
						d = Right
					}
					g.links[d] = s

					s.parent = f.parent
					f.parent = s

					g = s
					s = f.links[Left]
					s.parent = f
				}
				if (s.links[Left] == nil || s.links[Left].color == black) &&
					(s.links[Right] == nil || s.links[Right].color == black) {
					s.color = red
				} else {
					if s.links[Left] == nil || s.links[Left].color == black {
						y := s.links[Right]
						y.color = black
						s.color = red
						s.links[Right] = y.links[Left]
						y.links[Left] = s
						if s.links[Right] != nil {
							s.links[Right].parent = s
						}
						f.links[Left] = y
						s = y
						s.links[Left].parent = s
					}
					s.color = f.color
					f.color = black
					s.links[Left].color = black

					f.links[Left] = s.links[Right]
					s.links[Right] = f

					d := Left
					if g.links[Left] != f {
						d = Right
					}
					g.links[d] = s

					s.parent = f.parent
					f.parent = s
					if f.links[Left] != nil {
						f.links[Left].parent = f
					}
					break
				}
			}
			tmp = f
			f = f.parent
			if f == nil {
				f = (*prbnode)(unsafe.Pointer(&t.root))
			}
			d := Left
			if f.links[Left] != tmp {
				d = Right
			}
			dir = d
		}
	}
	w = nil
	t.count--

	return ret
}

func (t *PRbTree) Copy() *PRbTree {
	if t == nil {
		return nil
	}
	n := NewPRbTree(t.cmpFunc, t.extraParam)
	if n == nil {
		return nil
	}
	n.count = t.count
	if n.count == 0 {
		return n
	}
	var (
		x *prbnode
		y *prbnode
	)
	x = (*prbnode)(unsafe.Pointer(&t.root))
	y = (*prbnode)(unsafe.Pointer(&n.root))
	for {
		for x.links[Left] != nil {
			y.links[Left] = &prbnode{}
			y.links[Left].parent = y
			x = x.links[Left]
			y = y.links[Left]
		}
		y.links[Left] = nil
		for {
			y.data = x.data
			y.color = x.color
			if x.links[Right] != nil {
				y.links[Right] = &prbnode{}
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

func (t *PRbTree) Iter() Iterator {
	it := NewPRbIter()
	return it.HookWith(t)
}

type PRbIter struct {
	tree *PRbTree //the tree be iterated
	node *prbnode //current node in tree
}

func NewPRbIter() *PRbIter {
	return &PRbIter{}
}

func (it *PRbIter) HookWith(tree *PRbTree) *PRbIter {
	if it == nil {
		return nil
	}
	it.tree = tree
	it.node = nil
	return it
}

func (it *PRbIter) First() Item {
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

func (it *PRbIter) Last() Item {
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

func (it *PRbIter) Find(item Item) Item {
	if it == nil || it.tree == nil || item == nil {
		return nil
	}
	var (
		w *prbnode //walk node
		n *prbnode //child of w
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

func (it *PRbIter) Next() Item {
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

func (it *PRbIter) Prev() Item {
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

func (it *PRbIter) Current() Item {
	if it == nil || it.node == nil {
		return nil
	}
	return it.node.data
}

//don't change key part of item
func (it *PRbIter) Replace(new Item) Item {
	if it == nil || it.node == nil || new == nil {
		return nil
	}
	old := it.node.data
	it.node.data = new
	return old
}

func (it *PRbIter) CopyFrom(other *PRbIter) Item {
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

func (it *PRbIter) Insert(item Item) (*Item, bool) {
	if it == nil || it.tree == nil || item == nil {
		return nil, false
	}
	addr, ok := it.tree.insert(item)
	it.node = (*prbnode)(unsafe.Pointer(uintptr(unsafe.Pointer(addr)) - unsafe.Offsetof(it.node.data)))
	return addr, ok
}
