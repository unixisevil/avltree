package avltree

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"
)

const (
	insRandom = iota
	insAscending
	insDescending
	insBalanced
	insZigZag
	insAscendingShifted
	insCnt
)

const (
	delRandom = iota
	delReverse
	delSame
	delCnt
)

const (
	correctTest = iota
	overflowTest
)

func genBalancedTree(min, max int, ret []int) {
	if min > max {
		return
	}
	i := (min + max + 1) / 2
	ret[0] = i
	genBalancedTree(min, i-1, ret[1:len(ret)/2+1])
	genBalancedTree(i+1, max, ret[len(ret)/2+1:])
}

func genInsertArr(size int, order int) []int {
	arr := make([]int, size)
	switch order {
	case insRandom:
		for i := 0; i < size; i++ {
			arr[i] = i
		}
		rand.Shuffle(size, func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})

	case insAscending:
		for i := 0; i < size; i++ {
			arr[i] = i
		}
	case insDescending:
		for i := 0; i < size; i++ {
			arr[i] = size - 1 - i
		}
	case insBalanced:
		genBalancedTree(0, size-1, arr)
	case insZigZag:
		for i := 0; i < size; i++ {
			if i%2 == 0 {
				arr[i] = i / 2
			} else {
				arr[i] = size - 1 - i/2
			}
		}
	case insAscendingShifted:
		for i := 0; i < size; i++ {
			arr[i] = i + size/2
			if arr[i] >= size {
				arr[i] -= size
			}
		}
	}
	return arr
}

func genDeleteArr(insArr []int, order int) []int {
	arr := make([]int, len(insArr))
	switch order {
	case delRandom:
		for i := 0; i < len(insArr); i++ {
			arr[i] = i
		}
		rand.Shuffle(len(insArr), func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	case delReverse:
		for i := 0; i < len(insArr); i++ {
			arr[i] = insArr[len(insArr)-1-i]
		}
	case delSame:
		for i := 0; i < len(insArr); i++ {
			arr[i] = insArr[i]
		}
	}
	return arr
}

func (n *node) print(lvl int) {
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

func recurseVerifyTree(t *testing.T, node *node, ok *bool, count *int, min, max int, height *int) {
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
	recurseVerifyTree(t, node.links[Left], ok, &subcount[Left], min, d-1, &subheight[Left])
	recurseVerifyTree(t, node.links[Right], ok, &subcount[Right], d+1, max, &subheight[Right])

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
}

func verifyTree(t *testing.T, tree *AvlTree, arr []int) bool {
	ok := true
	if tree.Count() != len(arr) {
		t.Errorf("Tree count is %d, but should be %d.\n", tree.Count(), len(arr))
		ok = false
	}
	if ok {
		count := 0
		height := 0
		recurseVerifyTree(t, tree.root, &ok, &count, 0, math.MaxInt64, &height)
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
	return ok
}

func (t *AvlTree) print(title string) {
	fmt.Printf("%s: ", title)
	t.root.print(0)
	fmt.Println()
}

var intCmp Compare = func(a, b interface{}, extraParam interface{}) int {
	ia := a.(int)
	ib := b.(int)
	if ia < ib {
		return -1
	} else if ia > ib {
		return 1
	} else {
		return 0
	}
}

func (it *AvlIter) check(t *testing.T, i, n int, title string) bool {
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

func compareTrees(t *testing.T, a, b *node) bool {
	if a == nil && b == nil {
		return true
	}
	if a.data != b.data ||
		((a.links[Left] != nil) != (b.links[Left] != nil)) ||
		((a.links[Right] != nil) != (b.links[Right] != nil)) ||
		a.balance != b.balance {
		t.Logf("Copied nodes differ: a=%d (bal=%d) b=%d (bal=%d) a:", a.data, a.balance, b.data, b.balance)
		if a.links[Left] != nil {
			t.Logf("l")
		}
		if a.links[Right] != nil {
			t.Logf("r")
		}
		t.Logf(" b:")
		if b.links[Left] != nil {
			t.Logf("l")
		}
		if b.links[Right] != nil {
			t.Logf("r")
		}
		t.Log()
		return false
	}
	ok := true
	if a.links[Left] != nil {
		ok = ok && compareTrees(t, a.links[Left], b.links[Left])
	}
	if a.links[Right] != nil {
		ok = ok && compareTrees(t, a.links[Right], b.links[Right])
	}
	return ok
}

func testCorrectness(t *testing.T, insert, delete []int) (ok bool) {
	//测试创建树,插入数据
	tree := NewAvl(intCmp, nil)
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
		if !verifyTree(t, tree, insert[:i+1]) {
			ok = false
			return
		}
	}

	//测试修改树的同时使用迭代器访问树
	for i := 0; i < n; i++ {
		var (
			x AvlIter
			y AvlIter
			z AvlIter
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
		if !verifyTree(t, tree, insert) {
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
		if !verifyTree(t, tree, delete[i+1:]) {
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
			ok = ok && compareTrees(t, tree.root, copy.root)
		}

	}
	if ret := tree.Delete(insert[0]); ret != nil {
		t.Errorf("Deletion from empty tree succeeded.\n")
		ok = false
	}
	return
}

func iterFirst(t *testing.T, tree *AvlTree, n int) bool {
	var it AvlIter
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

func iterLast(t *testing.T, tree *AvlTree, n int) bool {
	var it AvlIter
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

func iterFind(t *testing.T, tree *AvlTree, n int) bool {
	var it AvlIter
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

func iterInsert(t *testing.T, tree *AvlTree, n int) bool {
	var it AvlIter
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

func iterNext(t *testing.T, tree *AvlTree, n int) bool {
	var it AvlIter
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

func iterPrev(t *testing.T, tree *AvlTree, n int) bool {
	var it AvlIter
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

func treeCopy(t *testing.T, tree *AvlTree, n int) bool {
	copy := tree.Copy()
	return compareTrees(t, tree.root, copy.root)
}

func testOverflow(t *testing.T, insert []int) bool {
	type testFunc func(t *testing.T, tree *AvlTree, n int) bool
	tests := [...]struct {
		name string
		fn   testFunc
	}{
		{"first item", iterFirst},
		{"last item", iterLast},
		{"find item", iterFind},
		{"insert item", iterInsert},
		{"next item", iterNext},
		{"prev item", iterPrev},
		{"copy tree", treeCopy},
	}
	n := len(insert)
	for _, test := range tests {
		if *verbose >= 2 {
			t.Logf("Running %s test...\n", test.name)
		}
		tree := NewAvl(intCmp, nil)
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
		if !verifyTree(t, tree, insert) {
			return false
		}
	}
	return true
}

func TestInt(t *testing.T) {
	switch *testMode {
	case correctTest:
		testCorrectness(t, insertArr, deleteArr)
	case overflowTest:
		testOverflow(t, insertArr)
	}
}

type kv struct {
	k string
	v int
}

type mkv struct {
	char rune
	pos  []int
}

var mapCmp Compare = func(a, b interface{}, extraParam interface{}) int {
	akv := a.(kv)
	bkv := b.(kv)
	if akv.k < bkv.k {
		return -1
	} else if akv.k > bkv.k {
		return 1
	} else {
		return 0
	}
}

var multiMapCmp Compare = func(a, b interface{}, extraParam interface{}) int {
	akv := a.(mkv)
	bkv := b.(mkv)
	if akv.char < bkv.char {
		return -1
	} else if akv.char > bkv.char {
		return 1
	} else {
		return 0
	}
}

func TestMap(t *testing.T) {
	m := NewAvl(mapCmp, nil)
	m.Insert(kv{"GPU", 15})
	m.Insert(kv{"RAM", 20})
	m.Insert(kv{"CPU", 10})

	it := m.Iter()
	for item := it.First(); item != nil; item = it.Next() {
		e := item.(kv)
		t.Logf("k = %v, v = %v\n", e.k, e.v)
	}
	t.Log()
	m.Replace(kv{"CPU", 25})
	m.Insert(kv{"SSD", 30})
	for item := it.First(); item != nil; item = it.Next() {
		e := item.(kv)
		t.Logf("k = %v, v = %v\n", e.k, e.v)
	}
}

func TestMultiMap(t *testing.T) {
	m := NewAvl(multiMapCmp, nil)
	str := "this is it"
	for pos, char := range str {
		if char == ' ' {
			continue
		}
		if item := m.Find(mkv{char: char}); item == nil {
			posArr := []int{pos}
			m.Replace(mkv{char: char, pos: posArr})
		} else {
			kv := item.(mkv)
			kv.pos = append(kv.pos, pos)
			m.Replace(kv)
		}
	}

	it := m.Iter()
	for item := it.First(); item != nil; item = it.Next() {
		e := item.(mkv)
		t.Logf("char = %c, pos = %v\n", e.char, e.pos)
	}
}

var treeSize = flag.Int("size", 15, "number of node in tree")
var testMode = flag.Int("mode", correctTest, "test mode of tree(0|1)")
var verbose = flag.Int("verbose", 0, "turn up test output message verbosity level(0|1|2|3)")
var insOrder = flag.Int("insOrder", insRandom, "insort array order(0|1|2|3|4|5)")
var delOrder = flag.Int("delOrder", delRandom, "delete array order(0|1|2)")

var (
	insertArr []int
	deleteArr []int
)

func TestMain(m *testing.M) {
	flag.Parse()
	rand.Seed(time.Now().Unix())

	if *insOrder < insRandom || *insOrder >= insCnt {
		fmt.Printf("invalid insertion order\n")
		os.Exit(1)
	}

	if *delOrder < delRandom || *delOrder >= delCnt {
		fmt.Printf("invalid delete order\n")
		os.Exit(1)
	}
	if *testMode != correctTest && *testMode != overflowTest {
		fmt.Printf("invalid test mode\n")
		os.Exit(1)
	}

	insertArr = genInsertArr(*treeSize, *insOrder)
	deleteArr = genDeleteArr(insertArr, *delOrder)
	if *verbose >= 1 {
		fmt.Printf("Insertion array: %v\n", insertArr)
		if *testMode == correctTest {
			fmt.Printf("Deletion array: %v\n", deleteArr)
		}
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}
