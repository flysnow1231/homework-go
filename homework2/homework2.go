package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
)

func main() {
	num := 10
	p := &num
	pointer(p)
	fmt.Printf("after update:%v\n", num)

	s := []int{2, 3, 4, 5, 6}
	doubleSlice(&s)
	fmt.Printf("double slice:%v\n", s)

	printNums()

	//runTask()

	r := Rectangle{3.11, 3.55}
	fmt.Printf("rectangle area is %v, Perimeter is %v \n", r.Area(), r.Perimeter())

	c := Circle{10}
	fmt.Printf("circle area is %v, Perimeter is %v \n", c.Area(), c.Perimeter())

	emp := Employee{
		Person: Person{
			Name: "Alice",
			Age:  30,
		},
		EmployeeID: 1001,
	}

	emp.PrintInfo()

	runChan()

	runMutex()

	atomicAdd()
}

func atomicAdd() {
	var (
		counter int64
		wg      sync.WaitGroup
	)

	// 启动 10 个 goroutine
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}

	wg.Wait()
	fmt.Println("counter by atomic is :", counter)
}

func runMutex() {
	var (
		counter int
		mu      sync.Mutex
		wg      sync.WaitGroup
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	fmt.Println("counter:", counter)
}

func pointer(p *int) {
	*p += 10
}

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person
	EmployeeID int
}

// PrintInfo 输出员工信息
func (e Employee) PrintInfo() {
	fmt.Printf(
		"EmployeeID: %d, Name: %s, Age: %d\n",
		e.EmployeeID,
		e.Name,
		e.Age,
	)
}

func doubleSlice(s *[]int) {

	for i := range *s {
		(*s)[i] *= 2
	}
}

func printNums() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 2; i <= 10; i += 2 {
			fmt.Printf("print even nums:%v\n", i)
		}
		defer wg.Done()
	}()
	wg.Add(1)
	go func() {
		for i := 1; i <= 9; i += 2 {
			fmt.Printf("print odd nums:%v\n", i)

		}
		defer wg.Done()
	}()
	wg.Wait()
}

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	height float64
	width  float64
}

type Circle struct {
	radius float64
}

func (r Rectangle) Area() float64 {
	return r.height * r.width
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.height + r.width)
}

func (r Circle) Perimeter() float64 {
	return math.Pi * r.radius * 2
}
func (r Circle) Area() float64 {
	return math.Pi * r.radius * r.radius
}

func producer(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 1; i <= 100; i++ {
		ch <- i
	}
	close(ch)
}

// 消费者：接收并打印
func consumer(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for v := range ch {
		fmt.Printf("consume: %v \n", v)
	}
}

func runChan() {
	// 创建一个带缓冲的通道（缓冲区大小 10）
	ch := make(chan int, 10)

	var wg sync.WaitGroup
	wg.Add(2)

	go producer(ch, &wg)
	go consumer(ch, &wg)

	wg.Wait()
	fmt.Println("done")
}
