package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LeonidS635/HyperLit/cmd"
	"github.com/LeonidS635/HyperLit/internal/app/hyperlit"
	"github.com/spf13/pflag"
)

const hlPath = "hl"

var path string

func init() {
	pflag.StringVar(&path, "path", "", "project path")
}

func main() {
	pflag.Parse()

	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)
	path = filepath.Join(exPath, path)

	hl := hyperlit.New(hlPath, path)
	defer hl.Clear()

	cmd.InitCmds(hl)
	if err := cmd.Execute(context.TODO()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
