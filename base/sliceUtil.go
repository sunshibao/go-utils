package base

import (
	"math/rand"
	"reflect"
	"time"
)

// 基于泛型的求切片的并集、交集、差集、删除、去重、判断存在
type iType interface {
	uint | uint16 | uint32 | uint64 | int | int16 | int32 | int64 | float32 | float64 | string | byte
}

type iTypeNum interface {
	uint | uint16 | uint32 | uint64 | int | int16 | int32 | int64 | float32 | float64
}

//SliceFindIndex 查找某个元素在切片的第一次出现位置Index
func SliceFindIndex[T iType](slice1 []T, value T) int {
	for i, e := range slice1 {
		if e == value {
			return i
		}
	}
	return -1
}

// SliceDel  切片删除指定元素 BY Index
func SliceDel[T iType](slice1 []T, index int) []T {
	if index >= len(slice1) {
		return slice1
	}
	// 将删除点前后的元素连接起来
	slice1 = append(slice1[:index], slice1[index+1:]...)

	return slice1
}

// SliceDelValue  切片删除指定元素 BY Value
func SliceDelValue[T iType](slice1 []T, value T) []T {
	//方法一：截取法（修改原切片）
	//这里利用对 slice 的截取删除指定元素。注意删除时，后面的元素会前移，所以下标 i 应该左移一位。
	//for i := 0; i < len(slice1); i++ {
	//	if slice1[i] == value {
	//		slice1 = append(slice1[:i], slice1[i+1:]...)
	//		i--
	//	}
	//}
	//return slice1

	//方法三：位移法
	//利用一个下标 index，记录下一个有效元素应该在的位置。遍历所有元素，当遇到有效元素，将其移动到 index 且 index 加一。最终 index 的位置就是所有有效元素的下一个位置，最后做一个截取就行了。这种方法会修改原来的 slice。
	//该方法可以看成对第一种方法截取法的改进，因为每次指需移动一个元素，性能更加。
	j := 0
	for _, v := range slice1 {
		if v != value {
			slice1[j] = v
			j++
		}
	}
	return slice1[:j]

	//方法四：位移法
	//创建了一个 slice，但是共用原始 slice 的底层数组。这样也不需要额外分配内存空间，直接在原 slice 上进行修改。
	//tgt := slice1[:0]
	//for _, v := range slice1 {
	//	if v != value {
	//		tgt = append(tgt, v)
	//	}
	//}
	//return tgt

	//从基准测试结果来看，性能最佳的方法是移位法，其中又属第一种实现方式较佳。性能最差的也是最常用的方法是截取法。随着切片长度的增加，上面四种删除方式的性能差异会愈加明显。
	//
	//实际使用时，我们可以根据不用场景来选择。如不能修改原切片使用拷贝法，可以修改原切片使用移位法中的第一种实现方式。
}

// SliceUnique  切片去重
func SliceUnique[T iType](slice1 []T) []T {
	//校验长度，如果长度为0或1时，直接返回原切片。即使长度为1时也不用去重。
	if len(slice1) < 2 {
		return slice1
	}
	newMap := make(map[T]int)
	newSlice := make([]T, len(newMap))
	for _, val := range slice1 {
		newMap[val]++
	}
	for _, val := range slice1 {
		_, ok := newMap[val]
		if ok {
			newSlice = append(newSlice, val)
			delete(newMap, val)
		}

	}
	return newSlice
}

// SliceIn  判断切片中是否含有某元素
func SliceIn[T iType](slice1 []T, value T) bool {
	for _, v := range slice1 {
		if v == value {
			return true
		}
	}

	return false
}

// SliceUnion  求并集
func SliceUnion[T iType](slice1, slice2 []T) []T {
	tempMap := make(map[T]int)
	for _, val := range slice1 {
		tempMap[val]++
	}
	for _, val := range slice2 {
		_, ok := tempMap[val]
		if !ok {
			slice1 = append(slice1, val)
		}
	}
	return slice1
}

// SliceIntersect 求切片的交集
func SliceIntersect[T iType](slice1, slice2 []T) []T {
	tempMap := make(map[T]int)
	newSlice := make([]T, 0)
	for _, val := range slice1 {
		tempMap[val]++
	}

	for _, val := range slice2 {
		_, ok := tempMap[val]
		if ok {
			newSlice = append(newSlice, val)
		}
	}
	return newSlice
}

// SliceDifference 求切片的差集 slice1-并集
func SliceDifference[T iType](slice1, slice2 []T) []T {
	tempMap := make(map[T]int)
	newSlice := make([]T, 0)
	inter := SliceIntersect(slice1, slice2)
	for _, val := range inter {
		tempMap[val]++
	}

	for _, value := range slice1 {
		_, ok := tempMap[value]
		if !ok {
			newSlice = append(newSlice, value)
		}
	}
	return newSlice
}

//SliceRand 切片乱序
func SliceRand[T iType](slice1 []T) []T {
	if len(slice1) < 2 {
		return slice1
	}
	swap := reflect.Swapper(slice1)
	rand.Seed(time.Now().Unix())
	for i := len(slice1) - 1; i >= 0; i-- {
		j := rand.Intn(len(slice1))
		swap(i, j)
	}
	return slice1
}

//Paginate 分页处理  1: 2: 3:   page第几页，size 一页多少个
func Paginate[T iType](x []T, page int, size int) []T {
	page = size * page
	limit := func() int {
		if page+size > len(x) {
			return len(x)
		} else {
			return page + size
		}

	}

	start := func() int {
		if page > len(x) {
			return len(x)
		} else {
			return page
		}

	}
	return x[start():limit()]
}

//SliceSort 标准冒泡排序，参数1为原数组，参数2如果为0则正序，为1则倒序。
func SliceSort[T iTypeNum](arr []T, sort int) []T {
	if len(arr) < 2 || sort > 1 {
		return arr //这里是直接返回
	}

	if sort == 0 {
		var temp T
		for i := 0; i < len(arr); i++ {
			for j := 0; j < len(arr)-(1+i); j++ {
				if arr[j] > arr[j+1] {
					temp = arr[j]
					arr[j] = arr[j+1]
					arr[j+1] = temp
					//arr[j],arr[j+1] = arr[j+1],arr[j]
				}
			}
		}
	} else if sort == 1 {
		for i := 0; i < len(arr); i++ {
			for j := 0; j < len(arr)-(1+i); j++ {
				if arr[j] < arr[j+1] {
					arr[j], arr[j+1] = arr[j+1], arr[j]
				}
			}
		}
	}

	return arr
}

// SliceSameElementShowNum 切片中相同元素的数量统计。返回：元素值与元素出现次数
func SliceSameElementShowNum[T iType](arr []T) map[T]int {
	//校验长度
	if len(arr) < 1 {
		return nil
	}
	//对象键值对法
	//该方法执行的速度比其他任何方法都快，就是占用的内存大一些
	tempMap := make(map[T]int, 0)

	//fmt.Println("初始，map的值：", tempMap, "，长度：", len(tempMap))

	for _, value := range arr {
		if _, ok := tempMap[value]; ok == true {
			tempMap[value]++
			//fmt.Println("是否有相同的key：", ok, value)
		} else {
			tempMap[value] = 1
		}
	}

	//fmt.Println("=======================遍历=======================")
	//for key, value := range tempMap {
	//	fmt.Println("map的key：", key, "------>value：", value)
	//}

	return tempMap
}

//检查是否在切片内
func IsInSlice(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
