package cmd

type Context struct {
	Debug bool
}

var Cli struct {
	Run BenchCommand `cmd:"" help:"Run benchmarks"`
}
