# Learn Concurrent Go


#### Amdahl's Law 
- non-parallel parts of the program becomes a bottleneck

#### Gustafson's Law
- we put the free computing to other use

#### Process and Thread
- Process is a program running in an operating sytem
- Thread a construct that runs within a process
- Process starts with a single thread aka "main thread"
- Process have their own resources like memory. If it crash, it won't affect other processes
- Since processes are their own "world", communicating with other processes is a lil bit challenging
- When we need to synchronize processes, we typically use os tools and other apps like files, db, etc.
- Threads are lightweight alternative for concurrency kind of like microprocess within a process
- Threads share resources thus making it more efficient but also dangerous
- Threads, on the other hand, don't share stack space with each other
- When a local variable is created within a function, that lives in the stack

#### How I understand Process and Thread
- Let's say you have a kitchen (process) and staffs are the threads
- Kitchen aims to serve food to customers
- Staffs have their own tasks to achieve a bigger goal
- Kitchen has its own resources like burners, knives, ingredients, etc
- Staffs share the resources in the kitchen 


#### Concurrency vs Parallelism
- if one core, threads take turns, concurrency
- if multicore, threads can run at the same time in each core, parallelism
