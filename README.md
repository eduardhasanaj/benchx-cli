# BENCHX CLI [![Go Report Card](https://goreportcard.com/badge/github.com/eduardhasanaj/benchx-cli)](https://goreportcard.com/report/github.com/eduardhasanaj/benchx-cli)

Benchx is a command line tool which provides a handy experience for generating benchmark graphs by examining the output of:

`
go test ./... -bench . -benchmem
`
The graphs are stored at `./benchx-graphs` by default. At this time this is not configurable.

## Usage
### Benchmark Naming Conventions
`benchx-cli` will examine the test names from the output of `go test ./... -bench . -benchmem`.
The naming schema of the benchmark is:
```
    1        2       3
__________________________
|        |        |      |
BenchmarkFibonacciRecusion
```

Three parts of the naming conventions are:
1. every benchmark must start with `Benchmark` keyword
2. the second part can act as a group identifier;
in such case all groups must be acknowledged by using `--groups` flag and should be separated using space
3. is the name of the benchmark case

### Comands
```
Flags:
  -h, --help    Show context-sensitive help.

Commands:
  run
    Run benchmarks
```
At this time, the `benchx-cli` has only one command.
This commands has a flag called groups which you can use to group your benchmarks according to the [Name Conventions](#benchmark-naming-conventions).

### Example
To demonstrate the functionality of `benchx-cli` a fibonacci benchmark is performed to compare the performance between loop vs recursion implementation.

Test code
```
func BenchmarkFibonacciRecursion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FibonacciRecursion(10)
	}
}

func BenchmarkFibonacciLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FibonacciLoop(10)
	}
}
```

The following command is used to run and generate benchmark graphs:
```
benchx-cli run --groups Fibonacci
```

The following graphs are generated:

![Iterations](example/benchx-graphs/fibonacci_iterations.png?raw=true "Iterations")
![Speed](example/benchx-graphs/fibonacci_ns_op.png?raw=true "Speed")
![Memory](example/benchx-graphs/fibonacci_b_op.png?raw=true "Memor")
![Allocations](example/benchx-graphs/fibonacci_allocs_op.png?raw=true "Allocs")