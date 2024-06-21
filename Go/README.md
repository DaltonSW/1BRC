# Go

## 1st Completion - 2m 30.66s

Naive implementation, timed with zsh `time` command. 
- 149.50s attributed to user, 5.78s attributed to system

- I used bufio.Scanner for file input
- I used 4 different maps to store count, sum, min, and max
- I used strconv.ParseFloat to get the float out of the string after splitting
- Waited until the end to divide the average
- Put all the names in a slice, sorted that, then looped over that to print everything out alphabetically
- Formatting was handled with Sprintf, formatting all the floats with `%.1f` to round them
