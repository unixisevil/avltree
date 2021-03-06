package bbst

import (
	"flag"
	"fmt"
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

const (
	avlNoParent = iota
	avlWithParent
	rbNoParent
	rbWithParent
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

func TestInt(t *testing.T) {
	switch *testMode {
	case correctTest:
		switch *treeType {
		case avlNoParent:
			testAvlCorrectness(t, insertArr, deleteArr)
		case avlWithParent:
			testPAvlCorrectness(t, insertArr, deleteArr)
		case rbNoParent:
			testRbCorrectness(t, insertArr, deleteArr)
		case rbWithParent:
			testPRbCorrectness(t, insertArr, deleteArr)
		}
	case overflowTest:
		switch *treeType {
		case avlNoParent:
			testAvlOverflow(t, insertArr)
		case avlWithParent:
			testPAvlOverflow(t, insertArr)
		case rbNoParent:
			testRbOverflow(t, insertArr)
		case rbWithParent:
			testPRbOverflow(t, insertArr)
		}
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
	var m SymTab
	switch *treeType {
	case avlNoParent:
		m = NewAvlTree(mapCmp, nil)
	case avlWithParent:
		m = NewPAvlTree(mapCmp, nil)
	case rbNoParent:
		m = NewRbTree(mapCmp, nil)
	case rbWithParent:
		m = NewPRbTree(mapCmp, nil)
	}
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
	var m SymTab
	switch *treeType {
	case avlNoParent:
		m = NewAvlTree(multiMapCmp, nil)
	case avlWithParent:
		m = NewPAvlTree(multiMapCmp, nil)
	case rbNoParent:
		m = NewRbTree(multiMapCmp, nil)
	case rbWithParent:
		m = NewPRbTree(multiMapCmp, nil)
	}
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
var treeType = flag.Int("type", avlNoParent, "test tree type, 0(avlNoParent), 1(avlWithParent), 2(rbNoParent), 3(rbWithParent)")
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
	if *treeType < avlNoParent || *treeType > rbWithParent {
		fmt.Printf("invalid tree type\n")
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
