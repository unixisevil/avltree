package avltree

import (
	"fmt"
	"math"
	"testing"
)

func (n *pnode) print(lvl int) {
	if n == nil {
		return
	}
	if lvl > 16 {
		fmt.Printf("[...]")
		return
	}
	fmt.Printf("%v", n.data)
	if n.links[Left] != nil || n.links[Right] != nil {
		fmt.Printf("(")
		n.links[Left].print(lvl + 1)
		if n.links[Right] != nil {
			fmt.Printf(",")
			n.links[Right].print(lvl + 1)
		}
		fmt.Printf(")")
	}
}

func (t *PAvlTree) print(title string) {
	fmt.Printf("%s: ", title)
	t.root.print(0)
	fmt.Println()
}

func (it *PAvlIter) check(t *testing.T, i, n int, title string) bool {
	ok := true
	prev := it.Prev()
	actual := 0
	expect := 0
	if prev != nil {
		actual = prev.(int)
	} else {
		actual = -1
	}
	if i == 0 {
		expect = -1
	} else {
		expect = i - 1
	}

	if (i == 0 && prev != nil) || (i > 0 && (prev == nil || prev != i-1)) {
		t.Errorf("%s iter ahead of %d, but should be ahead of %d.\n", title, actual, expect)
		ok = false
	}
	it.Next()
	cur := it.Current()
	if cur == nil || cur != i {
		actual := 0
		if cur != nil {
			actual = cur.(int)
		} else {
			actual = -1
		}
		t.Errorf("%s iter at %d, but should be at %d.\n", title, actual, i)
		ok = false
	}
	next := it.Next()
	if next != nil {
		actual = next.(int)
	} else {
		actual = -1
	}
	if i == n-1 {
		expect = -1
	} else {
		expect = i + 1
	}
	if (i == n-1 && next != nil) || (i != n-1 && (next == nil || next != i+1)) {
		t.Errorf("%s iter behind %d, but should be behind %d.\n", title, actual, expect)
		ok = false
	}
	it.Prev()
	return ok
}

func comparePAvlTrees(t *testing.T, a, b *pnode) bool {
	if a == nil && b == nil {
		return true
	}
	pf := func(n *pnode) Item {
		if n.parent != nil {
			return n.parent.data
		}
		return -1
	}
	cf := func(n *pnode, dir int) string {
		if n.links[dir] != nil {
			return "has"
		}
		return "no"
	}
	if a.data != b.data ||
		((a.links[Left] != nil) != (b.links[Left] != nil)) ||
		((a.links[Right] != nil) != (b.links[Right] != nil)) ||
		((a.parent != nil) != (b.parent != nil)) ||
		(a.parent != nil && b.parent != nil && a.parent.data != b.parent.data) ||
		a.balance != b.balance {

		t.Logf("Copied nodes differ:\n"+
			"a=%d, bal %d, parent %d, %s left child, %s right child\n"+
			"b=%d, bal %d, parent %d, %s left child, %s right child\n",
			a.data, a.balance, pf(a), cf(a, Left), cf(a, Right),
			b.data, b.balance, pf(b), cf(b, Left), cf(b, Right))
		return false
	}
	ok := true
	if a.links[Left] != nil {
		ok = ok && comparePAvlTrees(t, a.links[Left], b.links[Left])
	}
	if a.links[Right] != nil {
		ok = ok && comparePAvlTrees(t, a.links[Right], b.links[Right])
	}
	return ok
}

func recurseVerifyPTree(t *testing.T, node *pnode, ok *bool, count *int, min, max int, height *int) {
	var (
		d         int
		subcount  [ChildNum]int
		subheight [ChildNum]int
	)
	if node == nil {
		*count = 0
		*height = 0
		return
	}
	d = node.data.(int)
	if min > max {
		t.Errorf("Parents of node %d constrain it to empty range %d...%d.\n",
			d, min, max)
		*ok = false
	} else if d < min || d > max {
		t.Errorf("Node %d is not in range %d...%d implied by its parents.\n", d, min, max)
		*ok = false
	}
	recurseVerifyPTree(t, node.links[Left], ok, &subcount[Left], min, d-1, &subheight[Left])
	recurseVerifyPTree(t, node.links[Right], ok, &subcount[Right], d+1, max, &subheight[Right])

	*count = 1 + subcount[Left] + subcount[Right]
	maxHeight := 0
	if subheight[Left] > subheight[Right] {
		maxHeight = subheight[Left]
	} else {
		maxHeight = subheight[Right]
	}
	*height = maxHeight + 1

	if subheight[Right]-subheight[Left] != int(node.balance) {
		t.Errorf("Balance factor of node %d is %d, but should be %d.\n",
			d, node.balance, subheight[Right]-subheight[Left])

		*ok = false
	} else if node.balance < -1 || node.balance > 1 {
		t.Errorf("Balance factor of node %d is %d.\n", d, node.balance)
		*ok = false
	}
	for i := 0; i < ChildNum; i++ {
		if node.links[i] != nil && node.links[i].parent != node {
			var pdata Item
			if node.links[i].parent != nil {
				pdata = node.links[i].parent.data
			} else {
				pdata = -1
			}
			t.Errorf("Node %d has parent %d (should be %d).\n",
				node.links[i].data, pdata, d)
		}
	}
}

func verifyPTree(t *testing.T, tree *PAvlTree, arr []int) bool {
	ok := true
	if tree.Count() != len(arr) {
		t.Errorf("Tree count is %d, but should be %d.\n", tree.Count(), len(arr))
		ok = false
	}
	n := len(arr)
	if ok {
		count := 0
		height := 0
		recurseVerifyPTree(t, tree.root, &ok, &count, 0, math.MaxInt64, &height)
		if count != len(arr) {
			t.Errorf("Tree has %d nodes, but should have %d.\n", count, len(arr))
			ok = false
		}
	}
	if ok {
		for _, elem := range arr {
			if ret := tree.Find(elem); ret == nil {
				t.Errorf("Tree does not contain expected value %d.\n", elem)
				ok = false
			}
		}
	}
	if ok {
		var (
			it   PAvlIter
			item Item
			i    int
		)
		prev := -1
		for i, item = 0, it.HookWith(tree).First(); i < 2*n && item != nil; i, item = i+1, it.Next() {
			if item.(int) <= prev {
				t.Errorf("Tree out of order: %d follows %d in traversal\n", item, prev)
				ok = false
			}
			prev = item.(int)
		}
		if i != n {
			t.Errorf("Tree should have %d items, but has %d in traversal\n", n, i)
			ok = false
		}
	}
	if ok {
		var (
			it   PAvlIter
			item Item
			i    int
		)
		next := math.MaxInt64
		for i, item = 0, it.HookWith(tree).Last(); i < 2*n && item != nil; i, item = i+1, it.Prev() {
			if item.(int) >= next {
				t.Errorf("Tree out of order: %d precedes  %d in traversal\n", item, next)
				ok = false
			}
			next = item.(int)
		}
		if i != n {
			t.Errorf("Tree should have %d items, but has %d in traversal\n", n, i)
			ok = false
		}
	}
	if ok {
		init := tree.Iter()
		first := tree.Iter()
		last := tree.Iter()
		first.First()
		last.Last()
		if cur := init.Current(); cur != nil {
			t.Errorf("Inited iter should be nil, but is actually %d.\n", cur)
			ok = false
		}
		next := init.Next()
		if next != first.Current() {
			t.Errorf("Next after nil should be %d, but is actually %d.\n", first.Current(), next)
			ok = false
		}
		init.Prev()
		prev := init.Prev()
		if prev != last.Current() {
			t.Errorf("Prev before nil should be %d, but is actually %d.\n", last.Current(), prev)
			ok = false
		}
		init.Next()
	}
	return ok
}

func testPAvlCorrectness(t *testing.T, insert, delete []int) (ok bool) {
	//测试创建树,插入数据
	tree := NewPAvl(intCmp, nil)
	ok = true
	n := len(insert)

	for i := 0; i < n; i++ {
		if *verbose >= 2 {
			t.Logf("Inserting %d...\n", insert[i])
		}
		addr, _ := tree.insert(insert[i])
		if addr == nil {
			if *verbose >= 0 {
				t.Logf("Inserting invalid item")
			}
			return
		}
		if *addr != insert[i] {
			t.Logf("Inserting duplicate item ")
		}
		if *verbose >= 3 {
			tree.print("After insert")
		}
		if !verifyPTree(t, tree, insert[:i+1]) {
			ok = false
			return
		}
	}

	//测试修改树的同时使用迭代器访问树
	for i := 0; i < n; i++ {
		var (
			x PAvlIter
			y PAvlIter
			z PAvlIter
		)
		if insert[i] == delete[i] {
			continue
		}
		if *verbose >= 2 {
			t.Logf("Checking traversal from item %d...\n", insert[i])
		}
		if x.HookWith(tree).Find(insert[i]) == nil {
			t.Errorf("Can't find item %d in tree!\n", insert[i])
			continue
		}
		ok = ok && x.check(t, insert[i], len(insert), "Predeletion")

		if *verbose >= 3 {
			t.Logf("Deleting  %d...\n", delete[i])
		}
		delval := tree.Delete(delete[i])
		if delval == nil || delval != delete[i] {
			ok = false
			if delval == nil {
				t.Errorf("Not find item: %v\n", delete[i])
			} else {
				t.Errorf("Wrong node %d returned.\n", delval)
			}
		}
		y.CopyFrom(&x)
		if *verbose >= 3 {
			t.Logf("Re-inserting item %d.\n", delete[i])
		}
		if addr, _ := z.HookWith(tree).Insert(delete[i]); addr == nil {
			if *verbose >= 3 {
				t.Errorf("Re-inserting item %d failed.\n", delete[i])
			}
			ok = false
			return
		}

		ok = ok && x.check(t, insert[i], len(insert), "Postdeletion")
		ok = ok && y.check(t, insert[i], len(insert), "Copied")
		ok = ok && z.check(t, delete[i], len(delete), "Insertion")
		if !verifyPTree(t, tree, insert) {
			ok = false
			return
		}
	}

	//测试删除数据的同时，制造树的副本
	for i := 0; i < n; i++ {
		if *verbose >= 2 {
			t.Logf("Deleting  %d...\n", delete[i])
		}
		delval := tree.Delete(delete[i])
		if delval == nil || delval != delete[i] {
			ok = false
			if delval == nil {
				t.Errorf("Not find item: %v\n", delete[i])
			} else {
				t.Errorf("Wrong node %d returned.\n", delval)
			}
		}
		if *verbose >= 3 {
			tree.print("After delete")
		}
		if !verifyPTree(t, tree, delete[i+1:]) {
			ok = false
			return
		}
		if *verbose >= 2 {
			t.Logf("Copying tree and comparing...\n")
		}
		{
			copy := tree.Copy()
			if copy == nil {
				if *verbose >= 2 {
					t.Errorf("copy return nil")
				}
				ok = false
				return
			}
			ok = ok && comparePAvlTrees(t, tree.root, copy.root)
		}

	}
	if ret := tree.Delete(insert[0]); ret != nil {
		t.Errorf("Deletion from empty tree succeeded.\n")
		ok = false
	}
	return
}

func iterPAvlFirst(t *testing.T, tree *PAvlTree, n int) bool {
	var it PAvlIter
	if ret := it.HookWith(tree).First(); ret == nil || ret != 0 {
		actual := 0
		if ret != nil {
			actual = ret.(int)
		} else {
			actual = -1
		}
		t.Errorf("First item test failed: expected 0, got %d\n", actual)
		return false
	}
	return true
}

func iterPAvlLast(t *testing.T, tree *PAvlTree, n int) bool {
	var it PAvlIter
	if ret := it.HookWith(tree).Last(); ret == nil || ret != n-1 {
		actual := 0
		if ret != nil {
			actual = ret.(int)
		} else {
			actual = -1
		}
		t.Errorf("Last item test failed: expected %d, got %d\n", n-1, actual)
		return false
	}
	return true
}

func iterPAvlFind(t *testing.T, tree *PAvlTree, n int) bool {
	var it PAvlIter
	it.HookWith(tree)
	for i := 0; i < n; i++ {
		if ret := it.Find(i); ret == nil || ret != i {
			actual := 0
			if ret != nil {
				actual = ret.(int)
			} else {
				actual = -1
			}
			t.Errorf("Find item test failed: expected %d, got %d\n", i, actual)
			return false
		}
	}
	return true
}

func iterPAvlInsert(t *testing.T, tree *PAvlTree, n int) bool {
	var it PAvlIter
	it.HookWith(tree)
	for i := 0; i < n; i++ {
		if ret, succ := it.Insert(i); ret == nil || succ {
			actual := -2
			if ret != nil {
				actual = (*ret).(int)
			} else {
				actual = -1
			}
			t.Errorf("Insert item test failed: inserted dup  %d, got %d\n", i, actual)
			return false
		}
	}
	return true
}

func iterPAvlNext(t *testing.T, tree *PAvlTree, n int) bool {
	var it PAvlIter
	it.HookWith(tree)
	for i := 0; i < n; i++ {
		if ret := it.Next(); ret == nil || ret != i {
			actual := 0
			if ret != nil {
				actual = ret.(int)
			} else {
				actual = -1
			}
			t.Errorf("Next item test failed: expected %d, got %d\n", i, actual)
			return false
		}
	}
	return true
}

func iterPAvlPrev(t *testing.T, tree *PAvlTree, n int) bool {
	var it PAvlIter
	it.HookWith(tree)
	for i := n - 1; i >= 0; i-- {
		if ret := it.Prev(); ret == nil || ret != i {
			actual := 0
			if ret != nil {
				actual = ret.(int)
			} else {
				actual = -1
			}
			t.Errorf("Prev item test failed: expected %d, got %d\n", i, actual)
			return false
		}
	}
	return true
}

func copyPAvlTree(t *testing.T, tree *PAvlTree, n int) bool {
	copy := tree.Copy()
	return comparePAvlTrees(t, tree.root, copy.root)
}

func testPAvlOverflow(t *testing.T, insert []int) bool {
	type testFunc func(t *testing.T, tree *PAvlTree, n int) bool
	tests := [...]struct {
		name string
		fn   testFunc
	}{
		{"first item", iterPAvlFirst},
		{"last item", iterPAvlLast},
		{"find item", iterPAvlFind},
		{"insert item", iterPAvlInsert},
		{"next item", iterPAvlNext},
		{"prev item", iterPAvlPrev},
		{"copy tree", copyPAvlTree},
	}
	n := len(insert)
	for _, test := range tests {
		if *verbose >= 2 {
			t.Logf("Running %s test...\n", test.name)
		}
		tree := NewPAvl(intCmp, nil)
		for i := 0; i < n; i++ {
			addr, succ := tree.insert(insert[i])
			if addr == nil || !succ {
				if addr == nil && *verbose >= 0 {
					t.Errorf("invalid tree state")
				} else if !succ {
					t.Errorf("find duplicate data in tree")
				}
				return false
			}
		}
		if !test.fn(t, tree, n) {
			return false
		}
		if !verifyPTree(t, tree, insert) {
			return false
		}
	}
	return true
}
