package main

import (
	"fmt"
	"time"
	"github.com/bingfa/common"
	"strings"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

type user struct {
	Name string
}

func main() {
	//send 异步， recv同步
	//geneID()

	//send 异步， recv同步
	//showNotifies()

	//多channel合并
	//multiMerge()
	//send 异步， recv同步
	//simepleTry()

	//全同步，循环接收数据判断是否取完
	//method3()

	//发送接收对象，以及判断channel类型
	//estChanType()

	//处理channel超时

	//search by goroutine使用场景
	//  https://www.cnblogs.com/suoning/p/7237444.html

	//pub/sub mode
	pubsubMode();
}

func pubsubMode(){
	//core: data structer
	p := common.NewPublisher(100*time.Millisecond, 10)
	defer p.Close()
	//TODO
	all := p.Subscribe()
	golang := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}
		return false
	})
	p.Publish("hello, world!")
	p.Publish("hello, golang!")
	go func() {
		for msg := range all {
			fmt.Println("all:", msg)
		}
	} ()

	go func() {
		for msg := range golang {
			fmt.Println("golang:", msg)
		}
	} ()
	// 运行一定时间后退出
	time.Sleep(3 * time.Second)
}

func estChanType() {
	userChan := make(chan interface{}, 1)
	u := user{Name: "nick"}
	userChan <- &u
	close(userChan)

	u1 := <-userChan
	u2, ok := u1.(*user)
	if !ok {
		fmt.Println("cant not convert.")
		return
	}
	fmt.Println(u2)
}

// 生成自增的整数
func xrange() chan int{
	ch := make(chan int)
	go func() {
		for i := 0; ; i++ {
			fmt.Printf("send data %d, time: %s \n", i, time.Now().Sub(startTime))
			ch <- i  // 直到信道索要数据，才把i添加进信道
		}
	}()
	return ch
}

//自增生成器，当流程向信道索要数据时，goroutine才向channel添加数据到
func geneID() {
	generator := xrange()
	for i := 0; i < 1000; i++ {
		fmt.Printf("receive data %d, time: %s \n", <-generator, time.Now().Sub(startTime))
	}
}

func get_notification(user string) chan string{
	 // 此处可以查询数据库获取新消息等等..

	ch := make(chan string)
	go func() {
		ch <- fmt.Sprintf("Hi %s, welcome to weibo.com!", user)
	}()
	return ch
}

// 签名更明确
func getNotifies(user string) <- chan string{
	 // 此处可以查询数据库获取新消息等等..

	ch := make(chan string)
	go func() {
		fmt.Printf("send data %s, time: %s \n", user, time.Now().Sub(startTime))
		ch <- fmt.Sprintf("Hi %s, welcome to weibo.com!", user)
	}()
	return ch
}

//多个channel数据合并
func multiMerge() {
	result := common.FanIn(common.Branch(1), common.Branch(2), common.Branch(3))
	for i := 0; i < 3; i++ {
		fmt.Println(<-result)
	}
}

//场景：消息数据应该来自一个独立的服务，这个服务只负责 返回某个用户的新的消息提醒
func showNotifies() {
	//获取消息调用
	/*jack := get_notification("jack")
	//获得消息的返回
	fmt.Println(<-jack)*/

	ch := getNotifies("jack")
	for v := range ch {
		fmt.Printf("receive data %s, time: %s \n", v, time.Now().Sub(startTime))
		fmt.Printf(v)
		break //deadLock if not use break
	}

	//基于channel通信？？，避免deadLock

}

func simepleTry() {
	//method1 method2，接收值定义区别
	method1()
	//method2()
}

func method3() {
	intChan := make(chan int, 10)
	for i := 0; i < 10; i++ {
		intChan <- i
	}
	close(intChan)
	for {
		v, ok := <- intChan
		if !ok {
			fmt.Println("channel is close.")
			return
		}
		fmt.Println(v)
	}
}

func method2() {
	limitCnt := 10
	intChan := make(chan int, limitCnt)
	// func doSomething(top int) chan int
	intChan = doSomething(limitCnt)
	for i := 0;i <12 ;i++ {
		fmt.Println(<-intChan)
	}
}

func method1() {
	limitCnt := 10
	//func doSomething(top int) <-chan int / chan int
	intChan := doSomething(limitCnt)
	for v := range intChan {
		fmt.Println(v)
	}
	//取数据超出且发送数据的channel未关闭，会导致deadLock
	/*for i := 0;i <12 ;i++ {
		fmt.Println(<-intChan)
	}*/
}

func doSomething(top int) chan int {
	out := make(chan int, 10)
	go func() {
		for i := 0; i < top; i++ {
			//take too much time
			time.Sleep(500 * time.Millisecond)
			out <- i
		}
		//close channel in func
		close(out)
	}()
	return out
}
