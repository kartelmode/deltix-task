# deltix-task

There is a concurrent, but slow(lol) solution. I divide all rows from the user_data file by blocks, where in one block time is located in the one timestamp block. So then I procceed these blocks in goroutines. 

## Why I can divide rows into blocks?

min(total[user][t], l <= t < r) = total[user][l-1]+min(changed[user][t], l <= t < r)

The same with maximum.

And an average can be calculated as total[user][l-1]+sum(changed[user][t], l <= t < r)/deltaTime.