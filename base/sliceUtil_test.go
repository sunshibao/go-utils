package base

import (
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	"testing"
)

func TestSliceFindIndex(t *testing.T) {
	slice3 := []string{"aa", "bb", "cc", "dd", "8", "20", "8"}
	index := SliceFindIndex(slice3, "8")
	fmt.Println("slice3获取某个元素的下标：", index)
}

func TestSliceDelete(t *testing.T) {
	slice3 := []string{"aa", "bb", "cc", "dd", "8", "20", "8"}
	slice4 := SliceDel(slice3, 3)
	fmt.Println("slice3去掉下标为三的元素：", slice4)
}

func TestSliceUnique(t *testing.T) {
	slice3 := []string{"aa", "bb", "cc", "cc", "8", "20", "8"}
	unique := SliceUnique(slice3)
	fmt.Println("slice3去重为：", unique)
}

func TestSliceDelValue(t *testing.T) {
	slice3 := []string{"aa", "bb", "cc", "cc", "8", "20", "8"}
	value := SliceDelValue(slice3, "8")
	fmt.Println(fmt.Println("slice3去除某个元素：", value))
}

func TestSliceIn(t *testing.T) {
	slice3 := []string{"aa", "bb", "cc", "cc", "8", "20", "8"}
	isExist := SliceIn(slice3, "6")
	fmt.Println("slice3是否存在某元素：", isExist)
}

func TestSliceUnion(t *testing.T) {
	slice1 := []string{"1", "2", "7", "20", "8"}
	slice2 := []string{"2", "7", "5", "0"}
	un := SliceUnion[string](slice1, slice2)
	fmt.Println("slice1与slice2的并集为：", un)
}

func TestSliceIntersect(t *testing.T) {
	slice1 := []string{"1", "2", "7", "20", "8"}
	slice2 := []string{"2", "7", "5", "0"}
	in := SliceIntersect(slice1, slice2)
	fmt.Println("slice1与slice2的交集为：", in)
}

func TestSliceSliceDifference(t *testing.T) {
	slice1 := []string{"1", "2", "7", "20", "8"}
	slice2 := []string{"2", "7", "5", "0"}
	di := SliceDifference(slice1, slice2)
	fmt.Println("slice1与slice2的差集为：", di)
}

func TestSliceRand(t *testing.T) {
	slice3 := []string{"1", "2", "7", "20", "8"}
	for k := 0; k < 10000; k++ {
		slice3 = append(slice3, gconv.String(k))
	}
	di := SliceRand(slice3)
	fmt.Println("slice3 为切片打乱顺序：", di)
}

func TestPaginate(t *testing.T) {
	slice3 := []string{"1", "2", "7", "20", "8", "9", "dd", "22", "33", "44", "55"}
	di := Paginate(slice3, 1, 4)
	fmt.Println("slice3 分页展示：", di)
}

func TestSliceSort(t *testing.T) {
	slice3 := []int64{1, 2, 7, 20, 8, 9, 22, 33, 44, 55}
	di := SliceSort(slice3, 1)
	fmt.Println("slice3 倒序排列：", di)
}

func TestSliceSameElementShowNum(t *testing.T) {
	slice3 := []int64{1, 2, 7, 20, 8, 9, 8, 8, 8, 22, 1, 22, 33, 44, 55}
	di := SliceSameElementShowNum(slice3)
	fmt.Println("slice3 元素值出现的次数：", di)
}
