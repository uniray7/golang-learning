package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Res struct {
	Msg string
	Err error
}

func bye(ch chan<- Res, wg *sync.WaitGroup) {
	fmt.Println("wait...")
	wg.Wait()
	fmt.Println("finish waiting")

	close(ch)
}

func master() {
	num := 10

	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan Res, num)
	wg := &sync.WaitGroup{}

	defer bye(ch, wg)

	for i := 0; i < num; i++ {
		wg.Add(1)
		go worker(ctx, strconv.Itoa(i), ch, wg)
	}

	cnt := 0
	for msg := range ch {
		fmt.Println(msg)

		if msg.Err != nil {
			cancel()
			return
		}

		cnt++
		if cnt == num {
			fmt.Println("=======all finished")
			break
		}
	}

}

func main() {
	master()
	fmt.Println("byby")

	time.Sleep(11 * time.Second)

}

func worker(ctx context.Context, name string, ch chan<- Res, wg *sync.WaitGroup) {
	defer wg.Done()
	defer fmt.Println(name, " bye")

	if name == "2" || name == "3" {
		time.Sleep(1 * time.Second)
	}
	select {
	case <-ctx.Done():
		fmt.Println(name, " Cancel")
		return
	default:
		fmt.Printf("%s working\n", name)
		if name == "1" {
			ch <- Res{"", fmt.Errorf("gg")}
			//ch <- Res{fmt.Sprintf("%s finished", name), nil}

		} else {
			cnt := 0
			for i := 0; i < 1000; i++ {
				cnt += i
				if i%500 == 0 {
					time.Sleep(1 * time.Second)
				}
			}
			fmt.Println(name, cnt)
			ch <- Res{fmt.Sprintf("%s finished", name), nil}
		}
		return
	}
}
