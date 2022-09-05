package main

import (
	"github.com/alecthomas/kong"
	cmd "github.com/eduardhasanaj/benchx-cli/commands"
)

func main() {
	ctx := kong.Parse(&cmd.Cli)

	err := ctx.Run(&cmd.Context{})
	ctx.FatalIfErrorf(err)
}
