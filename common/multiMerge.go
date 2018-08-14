package common

import (
	"time"
	"math/rand"
)

//假如是比较耗时的计算
func DoStuff(x int) int {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	return 100 - x
}

//分支goroutine计算，并把计算结果流入各自信道
func Branch(x int) chan int{
	ch := make(chan int)
	go func() {
		ch <- DoStuff(x)
	}()
	return ch
}

//分支合并
func FanIn(chs... chan int) chan int {
	ch := make(chan int)

	for _, c := range chs {
		// 注意此处明确传值
		go func(c chan int) {
			ch <- <- c
		}(c) // 复合
	}

	return ch
}


/*func main() {
	result := FanIn(Branch(1), Branch(2), Branch(3))

	for i := 0; i < 3; i++ {
		fmt.Println(<-result)
	}
}*/

