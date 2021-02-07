package bbst

import (
	"fmt"
	"math"
	"testing"
)

func (n *prbnode) print(lvl int) {
	if n == nil {
		return
	}
	if lvl > 16 {
		fmt.Printf("[...]")
		return
	}
	fmt.Printf("%v[%d]", n.data, n.color)
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

func recurseVerifyPRbTree(t *testing.T, node *prbnode, ok *bool, count *int, min, max int, bh *int) {
	var (
		d        int           //data of tree node
		subcount [ChildNum]int //count of subtree
		subbh    [ChildNum]int //black height of subtree
	)
	if node == nil {
		*count = 0
		*bh = 0
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
	recurseVerifyPRbTree(t, node.links[Left], ok, &subcount[Left], min, d-1, &subbh[Left])
	recurseVerifyPRbTree(t, node.links[Right], ok, &subcount[Right], d+1, max, &subbh[Right])

	*count = 1 + subcount[Left] + subcount[Right]
	h := 0
	if node.color == black {
		h = 1
	}
	*bh = h + subbh[0]
	if node.color != red && node.color != black {
		t.Errorf("Node %d is neither red nor black (%d).\n", d, node.color)
		*ok = false
	}
	if node.color == red {
		if node.links[Left] != nil && node.links[Left].color == red {
			t.Errorf("Red node %d has red left child %d\n", d, node.links[Left].data)
			*ok = false
		}
		if node.links[Right] != nil && node.links[Right].color == red {
			t.Errorf("Red node %d has red right child %d\n", d, node.links[Right].data)
			*ok = false
		}
	}
	if subbh[Left] != subbh[Right] {
		t.Errorf("Node %d has two different black-heights: left bh=%d, right bh=%d\n", d, subbh[Left], subbh[Right])
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

			*ok = false
		}
	}
}

func verifyPRbTree(t *testing.T, tree *PRbTree, arr []int) bool {
	ok := true
	n := len(arr)
	if tree.Count() != n {
		t.Errorf("Tree count is %d, but should be %d.\n", tree.Count(), n)
		ok = false
	}
	if ok {
		if tree.root != nil && tree.root.color != black {
			t.Errorf("Tree root is not black.\n")
			ok = false
		}
	}
	if ok {
		count := 0
		bh := 0
		recurseVerifyPRbTree(t, tree.root, &ok, &count, 0, math.MaxInt64, &bh)
		if count != n {
			t.Errorf("Tree has %d nodes, but should have %d.\n", count, n)
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
			it   PRbIter
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
			it   PRbIter
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

func (t *PRbTree) print(title string) {
	fmt.Printf("%s: ", title)
	t.root.print(0)
	fmt.Println()
}

func (it *PRbIter) check(t *testing.T, i, n int, title string) bool {
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

func comparePRbTrees(t *testing.T, a, b *prbnode) bool {
	if a == nil && b == nil {
		return true
	}
	pf := func(n *prbnode) Item {
		if n.parent != nil {
			return n.parent.data
		}
		return -1
	}
	cf := func(n *prbnode, dir int) string {
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
		a.color != b.color {

		t.Logf("Copied nodes differ:\n"+
			"a=%d, color %d, parent %d, %s left child, %s right child\n"+
			"b=%d, color %d, parent %d, %s left child, %s right child\n",
			a.data, a.color, pf(a), cf(a, Left), cf(a, Right),
			b.data, b.color, pf(b), cf(b, Left), cf(b, Right))

		return false
	}
	ok := true
	if a.links[Left] != nil {
		ok = ok && comparePRbTrees(t, a.links[Left], b.links[Left])
	}
	if a.links[Right] != nil {
		ok = ok && comparePRbTrees(t, a.links[Right], b.links[Right])
	}
	return ok
}

func testPRbCorrectness(t *testing.T, insert, delete []int) (ok bool) {
	//测试创建树,插入数据
	tree := NewPRbTree(intCmp, nil)
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
		if !verifyPRbTree(t, tree, insert[:i+1]) {
			ok = false
			return
		}
	}

	//测试修改树的同时使用迭代器访问树
	for i := 0; i < n; i++ {
		var (
			x PRbIter
			y PRbIter
			z PRbIter
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
		if !verifyPRbTree(t, tree, insert) {
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
		if !verifyPRbTree(t, tree, delete[i+1:]) {
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
			ok = ok && comparePRbTrees(t, tree.root, copy.root)
		}
	}
	if ret := tree.Delete(insert[0]); ret != nil {
		t.Errorf("Deletion from empty tree succeeded.\n")
		ok = false
	}
	return
}

func prbIterFirst(t *testing.T, tree *PRbTree, n int) bool {
	var it PRbIter
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

func prbIterLast(t *testing.T, tree *PRbTree, n int) bool {
	var it PRbIter
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

func prbIterFind(t *testing.T, tree *PRbTree, n int) bool {
	var it PRbIter
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

func prbIterInsert(t *testing.T, tree *PRbTree, n int) bool {
	var it PRbIter
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

func prbIterNext(t *testing.T, tree *PRbTree, n int) bool {
	var it PRbIter
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

func prbIterPrev(t *testing.T, tree *PRbTree, n int) bool {
	var it PRbIter
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

func prbTreeCopy(t *testing.T, tree *PRbTree, n int) bool {
	copy := tree.Copy()
	return comparePRbTrees(t, tree.root, copy.root)
}

func testPRbOverflow(t *testing.T, insert []int) bool {
	type testFunc func(t *testing.T, tree *PRbTree, n int) bool
	tests := [...]struct {
		name string
		fn   testFunc
	}{
		{"first item", prbIterFirst},
		{"last item", prbIterLast},
		{"find item", prbIterFind},
		{"insert item", prbIterInsert},
		{"next item", prbIterNext},
		{"prev item", prbIterPrev},
		{"copy tree", prbTreeCopy},
	}
	n := len(insert)
	for _, test := range tests {
		if *verbose >= 2 {
			t.Logf("Running %s test...\n", test.name)
		}
		tree := NewPRbTree(intCmp, nil)
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
		if !verifyPRbTree(t, tree, insert) {
			return false
		}
	}
	return true
}
