package main

//这是一个用数组自己实现的队列，有助于对队列的理解，不建议用于生产中。
//切片实现固定长度队列，并加锁
import (
	"errors"
	"fmt"
	"log"
	"sync"
)

//使用一个结构体管理队列
type Queue struct {
	maxSize int64
	array   []int64 //数组=>模拟队列
	front   int64   //表示指向队列列首
	rear    int64   //表示指向队列的尾部
	flag    bool    //表示队尾是否在对头后面
	mut     *sync.Mutex
}

//实现结构体初始化有默认值
func NewQueue(max int64) Queue {
	return Queue{
		maxSize: max,
		array:   make([]int64, max),
		front:   0,
		rear:    -1,
		flag:    true,
		mut:     &sync.Mutex{},
	}
}

//添加数据到队列
func (this *Queue) AddQueue(val int64) (err error) {
	this.mut.Lock()
	defer this.mut.Unlock()
	//先判断队列是否已满
	if this.IsFull() {
		return errors.New("queue full")
	}
	this.rear++ //rear 后移
	if this.rear == this.maxSize {
		this.rear = 0
		this.flag = !this.flag
	}
	this.array[this.rear] = val
	return
}

//从队列中取出数据
func (this *Queue) GetQueue() (val int64, err error) {
	this.mut.Lock()
	defer this.mut.Unlock()
	//先判断队列是否空
	if this.IsEmpty() { //队空
		return -1, errors.New("queue empty")
	}
	val = this.array[this.front]
	this.front++
	if this.front == this.maxSize {
		this.front = 0
		this.flag = !this.flag
	}
	return val, err
}

//判断是否为空
func (this *Queue) IsEmpty() bool {
	if (this.flag == true && this.rear == this.front-1) || (this.flag == false && this.front == 0 && this.rear == this.maxSize-1) {
		return true
	} else {
		return false
	}
}

//判断是否为满
func (this *Queue) IsFull() bool {
	if (this.flag == true && this.front == 0 && this.rear == this.maxSize-1) || (this.flag == false && this.rear == this.front-1) {
		return true
	} else {
		return false
	}
}

//求队列中元素个数
func (this *Queue) QueueSize() int64 {
	if this.flag == false && this.front == 0 && this.rear == this.maxSize-1 {
		return 0
	}
	if this.flag == true {
		return this.rear - this.front + 1
	} else {
		return this.maxSize - (this.front - this.rear - 1)
	}
}

//显示队列,找到队首，然后到遍历到队尾
func (this *Queue) ShowQueue() {
	size := this.QueueSize()
	front := this.front
	for i := int64(1); i <= size; i++ {
		fmt.Printf("%d\t", this.array[front])
		front++
		if front == this.maxSize {
			front = 0
		}
	}
}

func main() {
	//先创建一个队列
	max := int64(50)
	queue := NewQueue(max)
	queue.ShowQueue()
	var wg sync.WaitGroup
	wg.Add(50)
	for i := 1; i <= 50; i++ {
		queue.AddQueue(int64(i))
	}
	queue.ShowQueue()
	for i := 1; i <= 50; i++ {
		go func() {
			out, err := queue.GetQueue()
			log.Println(err, queue.QueueSize(), "out:", out)
			wg.Done()
		}()
	}
	wg.Wait()
	queue.ShowQueue()
	fmt.Println(queue.QueueSize())
	fmt.Println(queue.GetQueue())
}
