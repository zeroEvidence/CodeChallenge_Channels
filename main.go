package main

import (
	"fmt"
	"sync"
)

func worker(jobs <-chan int32, results chan<- int, wg *sync.WaitGroup) {
	for {
		j, more := <-jobs

		if more {
			switch j % 3 {
			case 1:
				results <- int(j * 4)
			case 2:
				results <- int(j * 3)
			default:
			}
		} else {
			close(results)
			wg.Done()
			break
		}
	}
}

func sum(results <-chan int, sumRes chan<- int32, wg *sync.WaitGroup) {
	var sum int32

	for {
		r, more := <-results

		if more {
			sum += int32(r)
		} else {
			sumRes <- sum
			close(sumRes)
			wg.Done()
			break
		}
	}
}

func main() {
	var wg sync.WaitGroup
	jobs := make(chan int32)
	results := make(chan int)
	sumRes := make(chan int32, 1)

	wg.Add(1)
	go sum(results, sumRes, &wg)

	wg.Add(1)
	go worker(jobs, results, &wg)

	wg.Add(1)
	go func() {
		var i int32 = 1
		for i <= 1000000000 {
			if i%2 == 0 {
				i += 99
			}

			jobs <- i

			if i >= 1000000000 {
				close(jobs)
				wg.Done()
			}

			i++
		}
	}()

	wg.Wait()

	sum := <-sumRes
	fmt.Println(sum)
}
