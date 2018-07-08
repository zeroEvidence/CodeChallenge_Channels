# Problem

In simple point form write down what you see is wrong with the code. Please explain ​ WHY​ you
think it is wrong and how you would rectify it.

## Issues

- Worker function:
  - Line 7, 18:
    - What: Goroutine spawned for each int value that comes through the `jobs` channel.
    - Why: Creates unnessary overhead, overloads computer's resources. No performance gained over using the main process. Also not a process intensive task or a task that can benefit from sharing a thread as a coroutine.
    - How: Delete lines 7 and 18.

  - Lines 9, 10:
    - What: The expression in `case 0`.
    - Why: It asigns `j` to `j` as multiplying by 1 is frivioulous.
    - How: Delete both lines.

  - Lines 11 to 13:
    - What: j is assigned twice.
    - Why: Unnessasary extra operation.
    - How: Delete line 12 and rewrite line 13 to `results <- j * 4`.

  - Lines 14 to 17:
    - What: `j` is evaluated and passed to the results channel then reevaluated and assigned again.
    - Why: Opperation has no affect to the overall software, it is suppurvlous code that does nothing.
    - How: Delete line 16

- Main function:
  - Block from line 26 to 34:
    - Line 26, issue A:
      - What: For loop iterating `i` until `<= 1000000000`.
      - Why: `i` will be interated beyond its type boundry.
      - How: `i` should be declared outside the for loop, given the type of int32, subsequently, `i` will need to be iterated by +1 inside the loop, the `jobs` channel must have it's type change from int to int32, the parameter `jobs` in `worker` must be changed to reflact the change too, and the values given to the `results` channel to be typecast to an int by wrapping the values with `int()`

    - Line 26, issue B:
      - What: For loop iterating `i` until `<= 1000000000`.
      - Why: Large task that blocks the main thread.
      - How: I would have at least one Goroutine encompassing the entire `for` loop, which would get this process off the main thread thus allowing the main thread to be free for handling higher level software functions. Although, depending on the requirements, I would suggest that this task should be given to a dedicated pool of workers that autobalances the amount of workers available based on the amount of free resources available at runtime.

    - Line 27:
      - What: Goroutine spawned for each interation of the for loop.
      - Why: Creates unnessary overhead, overloads computer's resources, no performance gained over using a single Goroutine.
      - How: Delete lines 27 and 33, unless it was imperitive that `i` is to be evaluated and assigned asynchronously, and if it were, I would limit the amount of GoRoutines being used at anytime.

    - Lines 28 to 32: //Shared memory
      - What: `i` is being evaluated and assigned asynchronously relative to the outter for loop.
      - Why: Can produce an unpredictable amount and values of integers being passed into the jobs channel, because the time taken for a new Goroutine to be spawned is unknown, so the inner expressions will be evaluated in an unknown order and be unpredictable.
      - How: If this behaviour was not desirable then passing in `i` at the GoRoutine's execution and assigning it another variable name, perhaps `j` would make this more predictable because it passes in the value which existed at runtime when the Goroutine was called.

    - Line 35, issue A:
      - What: `jobs` channel is being closed before being read by `worker` function
      - Why: `worker` function cannot pull the values off the channel, and will end immediately.
      - How: defer the close against the channel or place the call at the end of main. However, please see issue "Line 35, issue B" for a better solution.

    - Line 35, issue B:
      - What: `jobs` channel is being closed outside of the closure sending values.
      - Why: The channel should be closed immediately after sending it's last value by the most inner closure possible, because it allows the rest of the software to synchronise expediantly. This is considered "best practice" too, as it also keeps the channel closures as close to it's use, making the code easy to read and understand.
      - How: add a clause inside the `for` loop to determine when the loop is finished sending data to the channel, i.e:
      ```
      if i <= 1000000000 {
        close(jobs)
      }
      ```

    - Lines 38 to 47:
      - What: `for` loop iterating `w` and `for` loop iterating `jobs2`.
      - Why: This code launches a 1000 Goroutines to launch the  `worker` function, which is unnessary because only one Goroutine is required to recieve values coming through on the channel, multiple goroutines waiting for values coming through channels will mean that only one goroutine will recieve a value, the others will miss out.
      - How: Delete lines 38 to 44 and 46 to 47, leaving only `go worker(w, jobs, results)`, delete the argument `w` and remove the `id` argument from the `worker` function.

    - Line 49:
      - What: `results` channel is being closed before being read by the `for` loop after it.
      - Why: `for` loop cannot pull the values off the channel, will immediately exit block of code.
      - How: The `results` channel should be closed by the function sending values once it's finished sending values. Please see issue at Line XX for a more indepth solution.

    - Line 51:
      - What: variable `sum` is being assigned its default empty value, `0`.
      - Why: unessary assignment.
      - How: delete the assigmnet to 0, i.e. ` = 0`.

    - Lines 53 to 55:
      - What: for loop is getting values from the `results` channel.
      - Why: the results channel is closed, this code block will immediately exit.
      - How: move lines 51 to 55 into its own function, begin function with a Goroutine before block on line 36 to 34.

