# Go

**Current Best: 36.669s**

## 1st implementation - 2m 30.66s (150.66)

Naive implementation, timed with zsh `time` command. 
- 149.50s attributed to user, 5.78s attributed to system

- I used bufio.Scanner for file input
- I used 4 different maps to store count, sum, min, and max
- I used strconv.ParseFloat to get the float out of the string after splitting
- Waited until the end to divide the average
- Put all the names in a slice, sorted that, then looped over that to print everything out alphabetically
- Formatting was handled with Sprintf, formatting all the floats with `%.1f` to round them

## 2nd implementation - 1m 20.785s (80.785)

Note: From here on, I'm going to be testing these using Go's built in benchmarking. I reran the first test using it, and it was basically the same as zsh`time`, so I don't think it matters to change the first one.

Wow what a change! I managed to cut the time in half before even implementing any coprocessing.

So what changed with this implementation?
- I learned more deeply about how structs worked. I wasn't confident enough in them to use them in my baseline implementation, so I just... didn't
  - I created a `city` struct that was comprised of a `count` int, and a `total`, `min`, and `max` float64
  - I then created a `mapHandler` struct that was just a map[string]city, to keep a mapping of city name to city info


*Nothing* changed as far as float parsing, calculations, or formatting is concerned, nor did I implement any goroutines. Plenty of optimization to be had!

## 3rd implementation - 3m 7.277s (187.277)

So... it got worse? Honestly not too surprised, this was my first attempt at implementing multithreading / go routines, and I've got no idea what the hell I'm doing!

But that's ok! Back to profiling and seeing what the bottleneck is now. 

## 4th implementation - 54.025s

Holy moly, down under a minute. Wild what happens when you actually research how to implement concurrency properly. There's definitely still some remnants of old mutex junk that can be cleaned up, and I'll probably futz around with some buffer size tweaking too to see where this can improve, but I'm pumped about the time savings. :)

## 5th implementation - 51.492s

lol I realized I had the const MBs but it wasn't multiplied, I'm a fool. Did this run wiht 32MB buffer for a slight improvement.

## 6th implementation - 36.669s

First attempt just had me increase the buffer size to 1000MB. Instantly got a crash, so reverted that. Decided to profile again, and saw a mess of concurrency lockin attempts.

Commented out the RLock/RUnlock for the Handler's city access, and that alone managed to drop the time down this much. While I don't have the expected min/max to compare, all of the averages were correct. I have to reasonably assume that it calculated min/max correctly if it was able to average correctly. (Hopefully that assumption won't be foiled by me having done something stupid...)
