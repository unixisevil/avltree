package bbst

const (
	Left = iota
	Right
	ChildNum
)

const (
	black = iota
	red
)

/*
  a  < b,  return negative value
  a  > b,  return positive value
  a == b,  return zero value
*/
type Compare func(a, b interface{}, extraParam interface{}) int

type Item interface{}

type Iterator interface {
	First() Item
	Last() Item
	Prev() Item
	Next() Item
	Current() Item
}

type SymTab interface {
	Count() int
	Find(target Item) Item
	Insert(item Item) bool
	Replace(item Item) Item
	Delete(item Item) Item
	Iter() Iterator
}
