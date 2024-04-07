package sort

import (
	"fmt"
	"testing"
)

// 选择排序
func TestSortSelect(*testing.T) {
	var aa = []int{1, 3, 2, 7, 4, 5}
	SortSelect(aa)
	fmt.Println(aa)
}

// 冒泡排序
func TestSortBubble(*testing.T) {
	var bb = []int{1, 3, 2, 7, 4, 5}
	SortBubble(bb)
	fmt.Println(bb)
}

// 插入排序
func TestSortInsert(*testing.T) {
	var bb = []int{1, 3, 2, 7, 4, 5}
	SortInsert(bb)
	fmt.Println(bb)
}

// 求前缀和数组
func Test(*testing.T) {
	var bb = []int{1, 3, 2, 7, 4, 5, -2}
	sum := PreSumArray(bb)

	l := 1
	r := 3
	fmt.Printf("数组范围l-r 累加和 sum:%d \n", GetSum(sum, l, r))

	l = 3
	r = 5
	fmt.Printf("数组范围l-r 累加和 sum:%d \n", GetSum(sum, l, r))
}
