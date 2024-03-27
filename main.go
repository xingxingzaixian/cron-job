package main

import (
	"cronJob/cmd"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	gctx.GetInitCtx()
	cmd.Execute()
}
