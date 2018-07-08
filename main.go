package main

import (
	"fmt"
	"sync"
)

func worker(id int, jobs <-chan int32, results chan<- int, wg *sync.WaitGroup) {
	for {
		j, more := <-jobs

		if more {
			// go func() {
			switch j % 3 {
			case 0:
				j = j * 1
			case 1:
				j = j * 2
				results <- int(j * 2)
			case 2:
				results <- int(j * 3)
				j = j * 3
			default:
			}
			// }()
		} else {
			close(results)
			wg.Done()
			return
		}
	}
}

func sum(results <-chan int, sumRes chan<- int32, wg *sync.WaitGroup) {
	var sum int32 = 0

	for {
		r, more := <-results

		if more {
			sum += int32(r)
		} else {
			sumRes <- sum
			close(sumRes)
			wg.Done()
			return
		}
	}
}

func main() {
	var wg sync.WaitGroup
	// i passed in will be outside of int range
	jobs := make(chan int32)
	results := make(chan int)
	sumRes := make(chan int32, 1)

	wg.Add(1)
	go sum(results, sumRes, &wg)

	// jobs2 := []int{}

	// for w := 1; w < 10; w++ {
	// 	jobs2 = append(jobs2, w)
	// }

	// for i, w := range jobs2 {
	wg.Add(1)
	go worker(10, jobs, results, &wg)
	// 	i = i + 1
	// }

	var i int32 = 1
	for i <= 1000000000 {
		if i%2 == 0 {
			i += 99
		}

		jobs <- i

		if i >= 1000000000 {
			close(jobs)
		}

		i++
	}

	wg.Wait()

	sum := <-sumRes
	fmt.Println(sum)
}
