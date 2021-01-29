# avltree

avl tree library implementation in golang

### Example

#### set:

```go
package main

import (
    "fmt"
    "math/rand"
    "time"

    "github.com/unixisevil/avltree"
)

func randRange(min, max int) int {
    return min + rand.Int()%(max-min+1)
}

var intCmp avltree.Compare = func(a, b interface{}, extraParam interface{}) int {
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

func main() {
    rand.Seed(time.Now().Unix())

    set := avltree.NewAvl(intCmp, nil)
    for i := 0; i < 10; i++ {
        set.Insert(randRange(100, 200))
    }

    it := set.Iter()
    for item := it.First(); item != nil; {
        num := item.(int)
        fmt.Printf("got: %d\n", num)
        //skip to next element before delete current element, avoid using invalid iterator
        item = it.Next()
        if num%2 == 0 {
            set.Delete(num)
        }
    }
    fmt.Println()
    for item := it.Last(); item != nil; item = it.Prev() {
        fmt.Printf("got: %d\n", item)
    }
}
```

#### map:

```go
package main

import (
    "fmt"

    "github.com/unixisevil/avltree"
)

type kv struct {
    k string
    v int
}

var mapCmp avltree.Compare = func(a, b interface{}, extraParam interface{}) int {
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

func main() {
    m := avltree.NewAvl(mapCmp, nil)
    m.Insert(kv{"GPU", 15})
    m.Insert(kv{"RAM", 20})
    m.Insert(kv{"CPU", 10})

    it := m.Iter()
    for item := it.First(); item != nil; item = it.Next() {
        e := item.(kv)
        fmt.Printf("k = %v, v = %v\n", e.k, e.v)
    }
    fmt.Println()
    m.Replace(kv{"CPU", 25})
    m.Insert(kv{"SSD", 30})
    for item := it.First(); item != nil; item = it.Next() {
        e := item.(kv)
        fmt.Printf("k = %v, v = %v\n", e.k, e.v)
    }

}
```

#### multi-map:

```go
package main

import (
	"fmt"

	"github.com/unixisevil/avltree"
)

type mkv struct {
	char rune
	pos  []int
}

var multiMapCmp avltree.Compare = func(a, b interface{}, extraParam interface{}) int {
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

func main() {
	m := avltree.NewAvl(multiMapCmp, nil)
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
		fmt.Printf("char = %c, pos = %v\n", e.char, e.pos)
	}
}
```

### running the tests:

```bash
go test -v  -run 'Int' -insOrder 0  -delOrder 0 -type 1 -size 1000  -verbose 3
```

```bash
go test -v  -run 'Map'  -mode 1
```
