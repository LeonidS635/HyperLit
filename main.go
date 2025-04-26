package main

import (
	"github.com/LeonidS635/HyperLit/cmd"
	"github.com/LeonidS635/HyperLit/internal/app/hyperlit"
)

const hlPath = "hl"

func main() {
	hl := hyperlit.New(hlPath, "testdata/toy_project")

	cmd.InitCmds(hl)
	cmd.Execute()
}
