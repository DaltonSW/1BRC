# Go

**Current Best: 80.785s**

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
