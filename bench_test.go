package bbst

import (
	"fmt"
	"testing"
)

func BenchmarkInsert(b *testing.B) {
	b.Run(fmt.Sprintf("avlNoParent/%d", *treeSize), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			tree := NewAvlTree(intCmp, nil)
			b.StartTimer()

			for _, elem := range insertArr {
				tree.Insert(elem)
			}
		}
	})
	b.Run(fmt.Sprintf("avlWithParent/%d", *treeSize), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			tree := NewPAvlTree(intCmp, nil)
			b.StartTimer()

			for _, elem := range insertArr {
				tree.Insert(elem)
			}
		}
	})
	b.Run(fmt.Sprintf("rbNoParent/%d", *treeSize), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			tree := NewRbTree(intCmp, nil)
			b.StartTimer()

			for _, elem := range insertArr {
				tree.Insert(elem)
			}
		}
	})
	b.Run(fmt.Sprintf("rbWithParent/%d", *treeSize), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			tree := NewPRbTree(intCmp, nil)
			b.StartTimer()

			for _, elem := range insertArr {
				tree.Insert(elem)
			}
		}
	})

}

func BenchmarkFind(b *testing.B) {
	b.Run(fmt.Sprintf("avlNoParent/%d", *treeSize), func(b *testing.B) {
		tree := NewAvlTree(intCmp, nil)
		for _, elem := range insertArr {
			tree.Insert(elem)
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, elem := range insertArr {
				tree.Find(elem)
			}
		}
	})
	b.Run(fmt.Sprintf("avlWithParent/%d", *treeSize), func(b *testing.B) {
		tree := NewPAvlTree(intCmp, nil)
		for _, elem := range insertArr {
			tree.Insert(elem)
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, elem := range insertArr {
				tree.Find(elem)
			}
		}
	})
	b.Run(fmt.Sprintf("rbNoParent/%d", *treeSize), func(b *testing.B) {
		tree := NewRbTree(intCmp, nil)
		for _, elem := range insertArr {
			tree.Insert(elem)
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, elem := range insertArr {
				tree.Find(elem)
			}
		}

	})
	b.Run(fmt.Sprintf("rbWithParent/%d", *treeSize), func(b *testing.B) {
		tree := NewPRbTree(intCmp, nil)
		for _, elem := range insertArr {
			tree.Insert(elem)
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, elem := range insertArr {
				tree.Find(elem)
			}
		}
	})

}

func BenchmarkDelete(b *testing.B) {
	b.Run(fmt.Sprintf("avlNoParent/%d", *treeSize), func(b *testing.B) {
		tree := NewAvlTree(intCmp, nil)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			for _, elem := range insertArr {
				tree.Insert(elem)
			}
			b.StartTimer()

			for _, elem := range insertArr {
				tree.Delete(elem)
			}
		}
	})
	b.Run(fmt.Sprintf("avlWithParent/%d", *treeSize), func(b *testing.B) {
		tree := NewPAvlTree(intCmp, nil)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			b.StopTimer()
			for _, elem := range insertArr {
				tree.Insert(elem)
			}
			b.StartTimer()

			for _, elem := range insertArr {
				tree.Delete(elem)
			}
		}
	})
	b.Run(fmt.Sprintf("rbNoParent/%d", *treeSize), func(b *testing.B) {
		tree := NewRbTree(intCmp, nil)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			b.StopTimer()
			for _, elem := range insertArr {
				tree.Insert(elem)
			}
			b.StartTimer()

			for _, elem := range insertArr {
				tree.Delete(elem)
			}
		}
	})
	b.Run(fmt.Sprintf("rbWithParent/%d", *treeSize), func(b *testing.B) {
		tree := NewPRbTree(intCmp, nil)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			b.StopTimer()
			for _, elem := range insertArr {
				tree.Insert(elem)
			}
			b.StartTimer()

			for _, elem := range insertArr {
				tree.Delete(elem)
			}
		}
	})
}
