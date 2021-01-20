package avltree

import (
	"unsafe"
)

const avlMaxHeight = 92

const (
	Left = iota
	Right
	ChildNum
)

/*
  a  < b,  return negative value
  a  > b,  return positive value
  a == b,  return zero value
*/
type Compare func(a, b interface{}, extraParam interface{}) int

type Item interface{}

type node struct {
	links   [ChildNum]*node
	data    Item
	balance int8
}

type AvlTree struct {
	root       *node       //root of  tree
	cmpFunc    Compare     //compare function
	extraParam interface{} //extra param for cmpFunc
	count      int         // number of item in tree
	generation int         // generation number
}

func NewAvl(cmp Compare, extra interface{}) *AvlTree {
	if cmp == nil {
		return nil
	}
	return &AvlTree{
		cmpFunc:    cmp,
		extraParam: extra,
	}
}

func (t *AvlTree) Count() int {
	if t == nil {
		return 0
	}
	return t.count
}

//search target in tree
//if find it return item
//else return nil
func (t *AvlTree) Find(target Item) Item {
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

func (t *AvlTree) insert(item Item) (*Item, bool) {
	if t == nil || item == nil {
		return nil, false
	}
	var (
		y   *node              //待更新平衡因子的最顶层节点
		z   *node              //y's  parent
		w   *node              //current walk node
		p   *node              //w's  parent
		n   *node              //new node
		r   *node              //new root node of rebalanced subtree
		dir byte               //下降方向
		da  [avlMaxHeight]byte //缓存的下降方向数组
		k   int                //length of da
	)
	z = (*node)(unsafe.Pointer(&t.root))
	dir = Left
	y = t.root
	for p, w = z, y; w != nil; p, w = w, w.links[dir] {
		cmp := t.cmpFunc(item, w.data, t.extraParam)
		if cmp == 0 {
			//fmt.Printf("item: %v, w.data: %v\n", item, w.data)
			return &w.data, false
		}
		if w.balance != 0 {
			z = p
			y = w
			k = 0
		}
		if cmp > 0 {
			dir = Right
		} else {
			dir = Left
		}
		da[k] = dir
		k++
	}
	n = &node{data: item}
	p.links[dir] = n
	t.count++
	if y == nil {
		//fmt.Println("tree is empty, ", n.data)
		return &n.data, true
	}
	for w, k = y, 0; w != n; w, k = w.links[da[k]], k+1 {
		if da[k] == Left {
			w.balance--
		} else {
			w.balance++
		}
	}
	if y.balance == -2 {
		x := y.links[Left]
		if x.balance == -1 {
			r = x
			y.links[Left] = x.links[Right]
			x.links[Right] = y
			x.balance = 0
			y.balance = 0
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
		}
	} else if y.balance == 2 {
		x := y.links[Right]
		if x.balance == 1 {
			r = x
			y.links[Right] = x.links[Left]
			x.links[Left] = y
			x.balance = 0
			y.balance = 0
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
		}
	} else {
		return &n.data, true
	}
	if y != z.links[Left] {
		dir = Right
	} else {
		dir = Left
	}
	z.links[dir] = r
	t.generation++
	return &n.data, true
}

//insert item in tree
//return true if item was successfully inserted
//return false if item already in tree
func (t *AvlTree) Insert(item Item) bool {
	_, succ := t.insert(item)
	return succ
}

//replace item in tree with same key item
func (t *AvlTree) Replace(item Item) Item {
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
func (t *AvlTree) Delete(item Item) Item {
	if t == nil || item == nil {
		return nil
	}

	var (
		pa  [avlMaxHeight]*node
		da  [avlMaxHeight]byte
		k   int
		w   *node
		dir byte
		cmp int
	)
	k = 0
	w = (*node)(unsafe.Pointer(&t.root))
	for cmp = -1; cmp != 0; cmp = t.cmpFunc(item, w.data, t.extraParam) {
		if cmp > 0 {
			dir = Right
		} else {
			dir = Left
		}
		pa[k] = w
		da[k] = dir
		k++
		w = w.links[dir]
		if w == nil {
			return nil
		}
	}
	ret := w.data

	//fmt.Printf("in delete(), ret: %v, k=%d\n", ret, k)
	if w.links[Right] == nil { //case 1, w has no right child
		pa[k-1].links[da[k-1]] = w.links[Left]
	} else { //case 2, w's right child has no left child
		r := w.links[Right]
		if r.links[Left] == nil {
			r.links[Left] = w.links[Left]
			r.balance = w.balance
			pa[k-1].links[da[k-1]] = r
			da[k] = Right
			pa[k] = r
			k++
		} else { //case 3, w's right child has left child

			var s *node
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
			s.links[Left] = w.links[Left]
			r.links[Left] = s.links[Right]
			s.links[Right] = w.links[Right]
			s.balance = w.balance

			pa[j-1].links[da[j-1]] = s
			da[j] = Right
			pa[j] = s
		}
	}
	w = nil
	//删除后，更新平衡因子, 重新平衡
	k--
	//fmt.Printf("before loop: k=%d\n", k)
	for ; k > 0; k-- {
		y := pa[k]
		if da[k] == Left {
			//fmt.Printf("left subtree branch, y: %v\n", *y)
			y.balance++
			if y.balance == 1 {
				break
			} else if y.balance == 2 { //重新平衡
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
					pa[k-1].links[da[k-1]] = r
				} else { /*  x.balance == 0  ||  x.balance == 1 */
					y.links[Right] = x.links[Left]
					x.links[Left] = y
					pa[k-1].links[da[k-1]] = x
					if x.balance == 0 {
						x.balance = -1
						y.balance = 1
						break
					} else {
						x.balance = 0
						y.balance = 0
					}
				}
			}
		} else {
			//fmt.Printf("else branch")
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
					pa[k-1].links[da[k-1]] = r
				} else {
					y.links[Left] = x.links[Right]
					x.links[Right] = y
					pa[k-1].links[da[k-1]] = x
					if x.balance == 0 {
						x.balance = 1
						y.balance = -1
						break
					} else {
						x.balance = 0
						y.balance = 0
					}
				}
			}
		}
	}

	t.count--
	t.generation++
	return ret
}

func (t *AvlTree) Copy() *AvlTree {
	if t == nil {
		return nil
	}
	n := NewAvl(t.cmpFunc, t.extraParam)
	if n == nil {
		return nil
	}
	n.count = t.count
	if n.count == 0 {
		return n
	}
	var (
		stack  [2 * (avlMaxHeight + 1)]*node
		height int
		x      *node
		y      *node
	)
	x = (*node)(unsafe.Pointer(&t.root))
	y = (*node)(unsafe.Pointer(&n.root))
	for {
		for x.links[Left] != nil {
			y.links[Left] = &node{}
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
			y.balance = x.balance
			if x.links[Right] != nil {
				y.links[Right] = &node{}
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

func (t *AvlTree) Iter() *AvlIter {
	it := NewIter()
	return it.HookWith(t)
}

type AvlIter struct {
	tree       *AvlTree            //the tree be iterated
	node       *node               //current node in tree
	stack      [avlMaxHeight]*node //all node above current node
	height     int                 //current depth of stack
	generation int                 // generation number
}

func NewIter() *AvlIter {
	return &AvlIter{}
}

func (it *AvlIter) HookWith(tree *AvlTree) *AvlIter {
	if it == nil {
		return nil
	}
	it.tree = tree
	it.node = nil
	it.height = 0
	it.generation = tree.generation

	return it
}

func (it *AvlIter) First() Item {
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

func (it *AvlIter) Last() Item {
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

func (it *AvlIter) Find(item Item) Item {
	if it == nil || it.tree == nil || item == nil {
		return nil
	}
	it.height = 0
	var (
		w *node //walk node
		n *node //child of w
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

func (it *AvlIter) Next() Item {
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

func (it *AvlIter) Prev() Item {
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

func (it *AvlIter) refresh() {
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

func (it *AvlIter) Current() Item {
	if it == nil || it.node == nil {
		return nil
	}
	return it.node.data
}

func (it *AvlIter) Replace(new Item) Item {
	if it == nil || it.node == nil || new == nil {
		return nil
	}
	old := it.node.data
	it.node.data = new
	return old
}

func (it *AvlIter) CopyFrom(other *AvlIter) Item {
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

func (it *AvlIter) Insert(item Item) (*Item, bool) {
	if it == nil || it.tree == nil || item == nil {
		return nil, false
	}
	addr, ok := it.tree.insert(item)

	it.node = (*node)(unsafe.Pointer(uintptr(unsafe.Pointer(addr)) - unsafe.Offsetof(it.node.data)))
	it.generation = it.tree.generation - 1
	return addr, ok
}
