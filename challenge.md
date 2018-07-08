# Problem

In simple point form write down what you see is wrong with the code. Please explain​ WHY​ you think it is wrong and how you would rectify it.

## Issues

- Worker function:

  - Lines 6 to 19:

    - What: This code block is the only place where values are being sent down the channel `results` and is not being closed.
    - Why: Closing a channel after a block is finished sending values is considered best practise.
    - How: Modify the `for` loop to an infinite loop, set variable `j, more` to be assigned `<-jobs` inside the for loop, encapsulate the `switch` with an `if` statement that evaluates the truthyness of `more`, `else` `close` `results` channel, call &WaitGroup.Done() (see issue labelled "Lines 22 to 58"), and break out of the `for` loop with a `return` or `break`; i.e:
      > ```
      >  func worker (..., wg *sync.WaitGroup) {
      >    for {
      >      j, more := <-jobs
      >
      >      if more {
      >        switch ...
      >      } else {
      >        close(results)
      >        wg.Done()
      >        break
      >      }
      >    }
      >  }
      > ```

  - Lines 7, 18:

    - What: Goroutine spawned for each value that comes through the `jobs` channel.
    - Why: Creates unnecessary overhead, overloads computer's resources. No performance gained over using the main process. Also not a process intensive task or a task that can benefit from sharing a thread as a coroutine.
    - How: Delete lines 7 and 18.

  - Lines 9, 10:

    - What: The expression in `case 0`.
    - Why: It assigns `j` to `j` as multiplying by 1 is frivolous.
    - How: Delete both lines.

  - Lines 11 to 13:

    - What: `j` is assigned twice.
    - Why: Unnecessary extra operation.
    - How: Delete line 12 and rewrite line 13 to: `results <- j * 4`.

  - Lines 14 to 17:

    - What: `j` is evaluated and passed to the results channel then is re-evaluated and assigned again.
    - Why: Operation has no affect to the overall software, it is superfluous code that does nothing.
    - How: Delete line 16

  - Lines 8 to 17:
    - What: `switch` statement.
    - Why: The switch statement has no default case, which is counter to what is considered best practice in this instance, as a value coming in from the channel j may get stuck if a suitable case is not found. This is handy to have in case of future modifications to the code, or if my recommended changes above is implemented, `case 0:` will no longer be an option and it will use the `default` case instead.
    - How: Add an empty `default` case after all the other cases.

- Main function:

  - Lines 22 to 58:

    - What: The `main` function has many responsibilities.
    - Why: The `main` function should have a single responsibility as per the S in the SOLID design principles.
    - How: Refactor `main` to have a single responsibility, (I haven't included this my solution due to not wanting to increase the complexity).

  - Lines 22 to 58:

    - What: The `main` function executes goroutines
    - Why: The `main` goroutine has no way of knowing when the other goroutines are done, so, (if I ignore the deadlocks), the `main` goroutine will exit before the other goroutines have finished.
    - How: Add a wait group from the `"sync"` package, declare it at the beginning of `main` before everything else, use the `&WaitGroup.Add(1)` method before a goroutine to add to the counter keeping track of the amount of goroutines, and have each goroutine call `&WaitGroup.Done()` when finishing, and apply `&WaitGroup.Wait()` before receiving and printing `sum`.

  - Lines 26 to 36:

    - What: The `jobs` channel is being sent values in this block.
    - Why: At this stage during the software's runtime, the software does not have any code that is awaiting values on the `jobs` channel, so all values being passed into it will be dropped.
    - How: This blocked must be moved to a point after the blocks of code which awaits on the `jobs` channel and subsequently the `results` channel, i.e. after Line 53.

  - Line 26:

    - What: `for` loop iterating `i` until `<= 1000000000`.
    - Why: `i` will be iterated beyond its type boundary.
    - How: `i` should be declared outside the `for` loop and given the type of `int32`, subsequently, `i` will need to be iterated by +1 inside the loop, the `jobs` channel must have it's type change from `int` to `int32`, the parameter `jobs` in the function `worker` must be changed to reflect the change too, and the values given to the `results` channel to be typecast to an `int` by wrapping the values with: `int()`

  - Lines 26 to 34:

    - What: `for` loop iterating `i` until `<= 1000000000`.
    - Why: Large task that blocks the main goroutine.
    - How: I would have at least one Goroutine encompassing the entire `for` loop block, which would allow this task to be executed asynchronously, thus allowing the main goroutine to be free for handling higher level software functions if needed. Although, depending on the requirements, I would suggest that this task should be given to a dedicated pool of workers that auto-balances the amount of workers available based on the amount of free resources available at runtime. This solution must use the `&WaitGroup.Add(1)` before execution and `&WaitGroup.Done()` when finishing, see above issues labelled "Lines 22 to 58", "Line 27" and "Lines 28 to 32" for more information.

  - Line 27:

    - What: Goroutine spawned for each iteration of the `for` loop.
    - Why: Creates unnecessary overhead, overloads computer's resources, no performance gained over using a single goroutine.
    - How: Delete lines 27 and 33, unless it was somehow imperative that `i` is to be evaluated and assigned asynchronously, and if it were, I would limit the amount of goroutines being used at any given time.

  - Lines 28 to 32:

    - What: `i` is being evaluated and assigned asynchronously relative to the outer for loop.
    - Why: This can produce an unpredictable amount and values of integers being passed into the jobs channel, because the time taken for a new Goroutine to be spawned is unknown, so the inner expressions will be evaluated in an unknown order and could be unpredictable.
    - How: (This solution depends on issue labelled "Line 27"), if this behaviour was not desirable then passing in the value `i` at runtime and assigning it another variable name, perhaps `j` would make this more predictable because it passes in the value which exists at runtime when the Goroutine was called.

  - Line 35, issue A:

    - What: `jobs` channel is being closed before being read by `worker` function
    - Why: `worker` function cannot pull the values off the channel, and will cause a deadlock.
    - How: defer the close of the channel or place the `close` call at the end of `main`. However, please see issue labelled "Line 35, issue B" for a better solution.

  - Line 35, issue B:

    - What: `jobs` channel is being closed outside of the closure sending values.
    - Why: The channel should be closed immediately after sending it's last value by the most inner closure possible, because it allows the rest of the software to synchronise expediently. This is considered "best practice" too, as it also keeps the channel closures as close to its use, making the code easy to read and understand.
    - How: Add a clause inside the `for` loop, at the bottom but before the iterator, to determine when the loop is finished sending data to the channel, and also tell the WaitGroup that we're done, (see issue labelled "Lines 22 to 58"). i.e:

    > ```
    >  ...
    >
    >  if i >= 1000000000 {
    >    close(jobs)
    >    wg.Done()
    >  }
    >
    >  i++
    > ```

  - Lines 38 to 47:

    - What: `for` loop iterating `w`, and `for` loop iterating `jobs2`.
    - Why: This code launches a 1000 goroutines and each execute the `worker` function, which is unnecessary because it creates excessive overheads, and only one goroutine is required to receive values coming through on the channel, so if there are multiple goroutines waiting for values coming through channels it means that only one goroutine will receive a value and the others will miss out causing deadlocks.
    - How: Delete lines 38 to 44 and 46 to 47, leaving only `go worker(w, jobs, results)`, delete the parameter `w` and remove the `id` parameter from the `worker` function, since there will only be one.

  - Line 49:

    - What: `close` of the `results` channel.
    - Why: The `for` loop directly after the `close` cannot pull the values off the channel, as a result, it will immediately exit the block of code.
    - How: The `results` channel should be closed by the function sending values once it has finished sending values. Please see the issue labelled "Lines 6 to 19" above for a more in depth solution.

  - Line 45:

    - What: Goroutine executes the worker function.
    - Why: The goroutine executes the worker function which waits for values coming through on the unbuffered channel `jobs`, this is a problem because by the time `worker` is executed, the `jobs` channel would've already dropped all the values that have been passed into it and is now closed; thus the goroutine is out of sync with the rest of the software and will create goroutines that that are forever awaiting values on the jobs channel, i.e. it creates a deadlock.
    - How: The goroutine must be moved before the block of code which is responsible for sending values down the `jobs` channel, i.e. before lines 26 to 34, and after the goroutine for calculating the sum from the `results` channel, (see issue labelled "Lines 53 to 55"); the `worker` must also `&WaitGroup.Add(1)` before execution and `&WaitGroup` to be passed into the `worker` function.

  - Line 51:

    - What: Variable `sum` is being assigned its default empty value, `0`.
    - Why: Unnecessary assignment.
    - How: Delete the assignment to 0, i.e. `= 0`.

  - Lines 53 to 55:
    - What: `for` loop is getting values from the `results` channel.
    - Why: The results channel is closed, and is out of sync with the rest of the software, this code block will create a deadlock.
    - How: Move lines 51 to 55 into its own function, add the required parameters, and back at `main` begin function with a goroutine before block on line 26 to 34 and before the `worker` goroutine, also add to the `&WaitGroup` with `&WaitGroup.Add()`. This also means we must add another int32 channel for the `sum` goroutine with a buffer of 1 (to prevent a deadlock), to pass back the summation of values being passed into the `results` channel. The summation must be given to the new channel after the closure of the `results` channel. To do so, we must apply the same code pattern as the solution for the issue labelled "Lines 6 to 19" above, but also passing the summation value to the new channel before closing it. The `main` function must wait, (using the &WaitGroup.Wait()), for the value before printing the value too. I.e:
      > ```
      >  func sum(results <-chan int, sumRes chan<- int32, wg *sync.WaitGroup) {
      >    var sum int32
      >
      >    for {
      >      r, more := <-results
      >
      >      if more {
      >        sum += int32(r)
      >      } else {
      >        sumRes <- sum
      >        close(sumRes)
      >        wg.Done()
      >        break
      >      }
      >    }
      >  }
      >
      >  ...
      >
      >  func main() {
      >    var wg sync.WaitGroup
      >    ...
      >    sumRes := make(chan int32, 1)
      >
      >    wg.Add(1)
      >    go sum(results, sumRes, &wg)
      >    ...
      >    wg.wait()
      >    sum := <-sumRes
      >    fmt.Println(sum)
      >  }
      > ```

The code should run without errors, and now look like:

> ```
>  package main
>
>  import (
>    "fmt"
>    "sync"
>  )
>
>  func worker(jobs <-chan int32, results chan<- int, wg *sync.WaitGroup) {
>    for {
>      j, more := <-jobs
>
>      if more {
>        switch j % 3 {
>        case 1:
>          results <- int(j * 4)
>        case 2:
>          results <- int(j * 3)
>        default:
>        }
>      } else {
>        close(results)
>        wg.Done()
>        break
>      }
>    }
>  }
>
>  func sum(results <-chan int, sumRes chan<- int32, wg *sync.WaitGroup) {
>    var sum int32
>
>    for {
>      r, more := <-results
>
>      if more {
>        sum += int32(r)
>      } else {
>        sumRes <- sum
>        close(sumRes)
>        wg.Done()
>        break
>      }
>    }
>  }
>
>  func main() {
>    var wg sync.WaitGroup
>    jobs := make(chan int32)
>    results := make(chan int)
>    sumRes := make(chan int32, 1)
>
>    wg.Add(1)
>    go sum(results, sumRes, &wg)
>
>    wg.Add(1)
>    go worker(jobs, results, &wg)
>
>    wg.Add(1)
>    go func() {
>      var i int32 = 1
>      for i <= 1000000000 {
>        if i%2 == 0 {
>          i += 99
>        }
>
>        jobs <- i
>
>        if i >= 1000000000 {
>          close(jobs)
>          wg.Done()
>        }
>
>        i++
>      }
>    }()
>
>    wg.Wait()
>
>    sum := <-sumRes
>    fmt.Println(sum)
>  }
> ```
