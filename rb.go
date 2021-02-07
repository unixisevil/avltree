package bbst

import (
	"unsafe"
)

const rbMaxHeight = 128

type rbnode struct {
	links [ChildNum]*rbnode //child node
	data  Item              //data item
	color byte              //node color
}

type RbTree struct {
	root       *rbnode     //root of  tree
	cmpFunc    Compare     //compare function
	extraParam interface{} //extra param for cmpFunc
	count      int         // number of item in tree
	generation int         // generation number
}

func NewRbTree(cmp Compare, extra interface{}) *RbTree {
	if cmp == nil {
		return nil
	}
	return &RbTree{
		cmpFunc:    cmp,
		extraParam: extra,
	}
}

func (t *RbTree) Count() int {
	if t == nil {
		return 0
	}
	return t.count
}

//search target in tree
//if find it return item
//else return nil
func (t *RbTree) Find(target Item) Item {
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

func (t *RbTree) insert(item Item) (*Item, bool) {
	if t == nil || item == nil {
		return nil, false
	}
	var (
		pa [rbMaxHeight]*rbnode //stack of rbnode
		da [rbMaxHeight]byte    //缓存的下降方向数组
		k  int                  //length of da
		w  *rbnode              //current walk node
		n  *rbnode              //new node
	)
	pa[0] = (*rbnode)(unsafe.Pointer(&t.root))
	da[0] = Left
	k = 1
	for w = t.root; w != nil; w = w.links[da[k-1]] {
		cmp := t.cmpFunc(item, w.data, t.extraParam)
		if cmp == 0 {
			return &w.data, false
		}
		pa[k] = w
		dir := Left
		if cmp > 0 {
			dir = Right
		}
		da[k] = byte(dir)
		k++
	}
	n = &rbnode{data: item, color: red}
	pa[k-1].links[da[k-1]] = n
	t.count++
	t.generation++
	for k >= 3 && pa[k-1].color == red {
		if da[k-2] == Left {
			/*
				   case 1, 新插入节点的n 的叔叔节点存在且是红色
				      pa[k-2](black)                       pa[k-2](red)
				      /      \                              /    \
				pa[k-1](red)   y(red)    =>       pa[k-1](black)   y(black)
				    /
				   n(red)
			*/
			y := pa[k-2].links[Right]
			if y != nil && y.color == red {
				pa[k-1].color = black
				y.color = black
				pa[k-2].color = red
				k -= 2
			} else {
				var x *rbnode
				/*
				 case 2, node n is left child of pa[k-1]
				 pa[k-2]|x (black)                      y(black)
				     /                                  /     \
				 pa[k-1]|y (red)      =>              n(red)  x(red)
				    /
				   n(red)
				*/
				if da[k-1] == Left {
					y = pa[k-1]
				} else {
					/*
					 case 3, node n is right child of pa[k-1], convert case 3 to case 2
					  pa[k-2](black)                  pa[k-2](black)
					    /                                 /
					 pa[k-1]|x(red)     =>                y(red)
					    \                               /
					    y|n (red)                      x(red)
					*/
					x = pa[k-1]
					y = x.links[Right]
					x.links[Right] = y.links[Left]
					y.links[Left] = x
					pa[k-2].links[Left] = y
				}
				x = pa[k-2]
				x.color = red
				y.color = black
				x.links[Left] = y.links[Right]
				y.links[Right] = x
				pa[k-3].links[da[k-3]] = y
				break
			}
		} else {
			y := pa[k-2].links[Left]
			if y != nil && y.color == red {
				pa[k-1].color = black
				y.color = black
				pa[k-2].color = red
				k -= 2
			} else {
				var x *rbnode
				if da[k-1] == Right {
					y = pa[k-1]
				} else {
					x = pa[k-1]
					y = x.links[Left]
					x.links[Left] = y.links[Right]
					y.links[Right] = x
					pa[k-2].links[Right] = y
				}
				x = pa[k-2]
				x.color = red
				y.color = black
				x.links[Right] = y.links[Left]
				y.links[Left] = x
				pa[k-3].links[da[k-3]] = y
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
func (t *RbTree) Insert(item Item) bool {
	_, succ := t.insert(item)
	return succ
}

//replace item in tree with same key item
func (t *RbTree) Replace(item Item) Item {
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
func (t *RbTree) Delete(item Item) Item {
	if t == nil || item == nil {
		return nil
	}
	var (
		pa  [rbMaxHeight]*rbnode //stack of rbnode
		da  [rbMaxHeight]byte    //缓存的下降方向数组
		k   int                  //length of da
		w   *rbnode              //current walk node
		cmp int
	)
	w = (*rbnode)(unsafe.Pointer(&t.root))
	for cmp = -1; cmp != 0; cmp = t.cmpFunc(item, w.data, t.extraParam) {
		dir := Left
		if cmp > 0 {
			dir = Right
		}
		pa[k] = w
		da[k] = byte(dir)
		k++
		w = w.links[dir]
		if w == nil {
			return nil
		}
	}
	ret := w.data
	if w.links[Right] == nil { //case 1, node to delete has no right child
		pa[k-1].links[da[k-1]] = w.links[Left]
	} else {
		r := w.links[Right]
		if r.links[Left] == nil { //case 2, node to delete w's right child has no left child
			r.links[Left] = w.links[Left]
			r.color, w.color = w.color, r.color //swap color
			pa[k-1].links[da[k-1]] = r          // hook w's right subtree with w's parent
			da[k] = Right
			pa[k] = r
			k++
		} else { //case 3, node to delete w's right child has left child
			var s *rbnode //w's successor
			j := k
			k++
			for {
				da[k] = Left
				pa[k] = r
				k++
				s = r.links[Left]
				if s.links[Left] == nil {
					break
				}
				r = s
			}
			//hook w's successor node s with w's parent
			da[j] = Right
			pa[j] = s
			pa[j-1].links[da[j-1]] = s

			//now r is s's parent node
			s.links[Left] = w.links[Left]
			r.links[Left] = s.links[Right]
			s.links[Right] = w.links[Right]
			s.color, w.color = w.color, s.color
		}
	}
	if w.color == black {
		for {
			x := pa[k-1].links[da[k-1]]
			if x != nil && x.color == red {
				x.color = black
				break
			}
			if k < 2 {
				break
			}
			if da[k-1] == Left {
				//node x's sibling
				s := pa[k-1].links[Right]
				if s.color == red {
					s.color = black
					pa[k-1].color = red
					pa[k-1].links[Right] = s.links[Left]
					s.links[Left] = pa[k-1]
					pa[k-2].links[da[k-2]] = s
					pa[k] = pa[k-1]
					da[k] = Left
					pa[k-1] = s
					k++
					s = pa[k-1].links[Right]
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
						pa[k-1].links[Right] = y
						s = y
					}
					s.color = pa[k-1].color
					pa[k-1].color = black
					s.links[Right].color = black

					pa[k-1].links[Right] = s.links[Left]
					s.links[Left] = pa[k-1]
					pa[k-2].links[da[k-2]] = s
					break
				}
			} else {
				//node x's sibling
				s := pa[k-1].links[Left]
				if s.color == red {
					s.color = black
					pa[k-1].color = red
					pa[k-1].links[Left] = s.links[Right]
					s.links[Right] = pa[k-1]
					pa[k-2].links[da[k-2]] = s
					pa[k] = pa[k-1]
					da[k] = Right
					pa[k-1] = s
					k++
					s = pa[k-1].links[Left]
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
						pa[k-1].links[Left] = y
						s = y
					}
					s.color = pa[k-1].color
					pa[k-1].color = black
					s.links[Left].color = black

					pa[k-1].links[Left] = s.links[Right]
					s.links[Right] = pa[k-1]
					pa[k-2].links[da[k-2]] = s
					break
				}
			}
			k--
		}
	}
	w = nil
	t.count--
	t.generation++
	return ret
}

func (t *RbTree) Copy() *RbTree {
	if t == nil {
		return nil
	}
	n := NewRbTree(t.cmpFunc, t.extraParam)
	if n == nil {
		return nil
	}
	n.count = t.count
	if n.count == 0 {
		return n
	}
	var (
		stack  [2 * (rbMaxHeight + 1)]*rbnode
		height int
		x      *rbnode
		y      *rbnode
	)
	x = (*rbnode)(unsafe.Pointer(&t.root))
	y = (*rbnode)(unsafe.Pointer(&n.root))
	for {
		for x.links[Left] != nil {
			y.links[Left] = &rbnode{}
			stack[height] = x
			height++
			stack[height] = y
			height++
			x = x.links[Left]
			y = y.links[Left]
		}
		y.links[Left] = nil
		for {
			y.data = x.data
			y.color = x.color
			if x.links[Right] != nil {
				y.links[Right] = &rbnode{}
				x = x.links[Right]
				y = y.links[Right]
				break
			} else {
				y.links[Right] = nil
			}
			if height <= 2 {
				return n
			}
			height--
			y = stack[height]
			height--
			x = stack[height]
		}
	}
}

func (t *RbTree) Iter() Iterator {
	it := NewRbIter()
	return it.HookWith(t)
}

type RbIter struct {
	tree       *RbTree              //the tree be iterated
	node       *rbnode              //current node in tree
	stack      [rbMaxHeight]*rbnode //all node above current node
	height     int                  //current depth of stack
	generation int                  // generation number
}

func NewRbIter() *RbIter {
	return &RbIter{}
}

func (it *RbIter) HookWith(tree *RbTree) *RbIter {
	if it == nil {
		return nil
	}
	it.tree = tree
	it.node = nil
	it.height = 0
	it.generation = tree.generation

	return it
}

func (it *RbIter) First() Item {
	if it == nil || it.tree == nil {
		return nil
	}
	it.height = 0
	w := it.tree.root
	if w == nil {
		return nil
	}
	for w.links[Left] != nil {
		it.stack[it.height] = w
		it.height++
		w = w.links[Left]
	}
	it.node = w
	return w.data
}

func (it *RbIter) Last() Item {
	if it == nil || it.tree == nil {
		return nil
	}
	it.height = 0
	w := it.tree.root
	if w == nil {
		return nil
	}
	for w.links[Right] != nil {
		it.stack[it.height] = w
		it.height++
		w = w.links[Right]
	}
	it.node = w
	return w.data
}

func (it *RbIter) Find(item Item) Item {
	if it == nil || it.tree == nil || item == nil {
		return nil
	}
	it.height = 0
	var (
		w *rbnode //walk node
		n *rbnode //child of w
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
		it.stack[it.height] = w
		it.height++
	}
	it.height = 0
	it.node = nil
	return nil
}

func (it *RbIter) Next() Item {
	if it == nil || it.tree == nil {
		return nil
	}
	if it.generation != it.tree.generation {
		it.refresh()
	}
	w := it.node
	if w == nil {
		return it.First()
	} else if w.links[Right] != nil {
		it.stack[it.height] = w
		it.height++
		w = w.links[Right]
		for w.links[Left] != nil {
			it.stack[it.height] = w
			it.height++
			w = w.links[Left]
		}
	} else {
		for {
			if it.height == 0 {
				it.node = nil
				return nil
			}
			n := w
			it.height--
			w = it.stack[it.height]
			if w.links[Right] != n {
				break
			}
		}
	}
	it.node = w
	return w.data
}

func (it *RbIter) Prev() Item {
	if it == nil || it.tree == nil {
		return nil
	}
	if it.generation != it.tree.generation {
		it.refresh()
	}
	w := it.node
	if w == nil {
		return it.Last()
	} else if w.links[Left] != nil {
		it.stack[it.height] = w
		it.height++
		w = w.links[Left]
		for w.links[Right] != nil {
			it.stack[it.height] = w
			it.height++
			w = w.links[Right]
		}
	} else {
		for {
			if it.height == 0 {
				it.node = nil
				return nil
			}
			n := w
			it.height--
			w = it.stack[it.height]
			if w.links[Left] != n {
				break
			}
		}

	}
	it.node = w
	return w.data
}

func (it *RbIter) refresh() {
	if it == nil || it.tree == nil {
		return
	}
	it.generation = it.tree.generation
	if it.node != nil {
		cmpFunc := it.tree.cmpFunc
		param := it.tree.extraParam
		node := it.node
		it.height = 0
		for w := it.tree.root; w != node; {
			it.stack[it.height] = w
			it.height++
			ret := cmpFunc(node.data, w.data, param)
			if ret > 0 {
				w = w.links[Right]
			} else {
				w = w.links[Left]
			}
		}
	}
}

func (it *RbIter) Current() Item {
	if it == nil || it.node == nil {
		return nil
	}
	return it.node.data
}

//don't change key part of item
func (it *RbIter) Replace(new Item) Item {
	if it == nil || it.node == nil || new == nil {
		return nil
	}
	old := it.node.data
	it.node.data = new
	return old
}

func (it *RbIter) CopyFrom(other *RbIter) Item {
	if it == nil || other == nil {
		return nil
	}
	if it != other {
		it.tree = other.tree
		it.node = other.node
		it.generation = other.generation
		if it.generation == it.tree.generation {
			it.height = other.height
			copy(it.stack[:it.height], other.stack[:other.height])
		}
	}
	if it.node == nil {
		return nil
	}
	return it.node.data
}

func (it *RbIter) Insert(item Item) (*Item, bool) {
	if it == nil || it.tree == nil || item == nil {
		return nil, false
	}
	addr, ok := it.tree.insert(item)

	it.node = (*rbnode)(unsafe.Pointer(uintptr(unsafe.Pointer(addr)) - unsafe.Offsetof(it.node.data)))
	it.generation = it.tree.generation - 1
	return addr, ok
}
