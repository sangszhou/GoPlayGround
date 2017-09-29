package main
import "sort"
import "fmt"

type ByLength []string

//We implement sort.Interface - Len, Less, and Swap -
// on our type so we can use the sort package’s generic Sort function.
// Len and Swap will usually be similar across types and Less will
// hold the actual custom sorting logic. In our case we want to sort
// in order of increasing string length, so we use len(s[i]) and len(s[j]) here.

// 这个例子简直不能再奇怪了，完全看不懂
// 真的很讨厌这个写法，golang 惹人讨厌的源泉之一
func (s ByLength) Len() int {
	return len(s)
}
func (s ByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

func main() {
	fruits := []string{"peach", "banana", "kiwi"}
	sort.Sort(ByLength(fruits))
	fmt.Println(fruits)
}
