package sort

// 选择排序O(N^2)
// arr[0~N-1]范围上，找到最小值所在的位置，然后把最小值交换到0位置。
// arr[1~N-1]范围上，找到最小值所在的位置，然后把最小值交换到1位置。
// arr[2~N-1]范围上，找到最小值所在的位置，然后把最小值交换到2位置，
// arr[N-2~N-1]范围上，找到最小值位置，然后把最小值交换到N-2位置
// 选择排序的时间复杂度为O(N^2)
func SortSelect(arr []int) {
	n := len(arr)
	if n < 2 {
		return
	}
	for i := 0; i < n-1; i++ {
		// i~n-1位置上找最小值
		// i位置的数和最小值的值交换

		minIndex := i

		for j := i + 1; j < n; j++ {
			if arr[j] < arr[minIndex] { //循环整个数组跟第一个比，如果更小把下标给第一个，以此类推
				minIndex = j
			}
		}
		swap(arr, i, minIndex) //全部比较完以后，找到最小值的下标，把值交换
	}
}

// 冒泡排序O(N^2)
// 在arr[0~N-1]范围上:
// arr[0]和arr[1]，谁大谁来到1位置;
// arr[1]和arr[2]，谁大谁来到2位置;
// arr[N-2]和arr[N-1]，谁大谁来到N-1位置;
// 在arr[0~N-2]范围上，重复上面的过程，最后一步是arr[N-3]和arr[N-2]，谁大谁来到N-2位置在arr[0~N-3]范围上，重复上面的过程，最后一步是arr[N-4]和arr[N-3]，谁大谁来到N-3位置
// 最后在arr[0~1]范围上，重复上面的过程，但最后一步是arr[0]和arr[1]，谁大谁来到1位置
func SortBubble(arr []int) {
	n := len(arr)
	if n < 2 {
		return
	}
	for end := n - 1; end > 0; end-- {
		for i := 0; i < end; i++ {
			if arr[i] > arr[i+1] {
				swap(arr, i, i+1)
			}
		}
	}
}

// 插入排序O(N^2)
// 想让arr[0~1]上有序，所以从arr[1]开始往前看，如果arr[1]<arr[0]，就交换。否则什么也不做。
// 想让arr[0~i]上有序，所以从arr[i]开始往前看，arr[i]这个数不停向左移动，一直移动到左边的数字不再比自己大，停止移动。最后一步，想让arr[0~N-1]上有序， arr[N-1]这个数不停向左移动，一直移动到左边的数字不再比自己大停止移动。
func SortInsert(arr []int) {
	n := len(arr)
	if n < 2 {
		return
	}
	for i := 0; i < n; i++ {
		for end := i; end > 0; end-- {
			if arr[i] < arr[i-1] {
				swap(arr, i, i-1)
			}
		}
	}
}

func swap(arr []int, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

// 前缀和数组 返回0-1，,0-2，0-3，0-4所有的和
func PreSumArray(arr []int) []int {
	n := len(arr)
	sum := make([]int, n)
	sum[0] = arr[0]
	for i := 1; i < n; i++ {
		sum[i] = sum[i-1] + arr[i]
	}
	return sum
}

// 利用前缀和 加工得到，比如求3~5 ，其实可以用0-5的和减去0-2的和。
func GetSum(sum []int, l, r int) int {
	if l == 0 {
		return sum[r]
	} else {
		return sum[r] - sum[l-1]
	}
}
