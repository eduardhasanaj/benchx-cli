package main

import (
	"github.com/alecthomas/kong"
	"github.com/eduardhasanaj/benchx-cli/cmd"
)

func main() {
	ctx := kong.Parse(&cmd.Cli)

	err := ctx.Run(&cmd.Context{})
	ctx.FatalIfErrorf(err)
}
