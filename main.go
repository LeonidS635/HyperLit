package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/LeonidS635/HyperLit/cmd"
	"github.com/LeonidS635/HyperLit/internal/app/hyperlit"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	hl := hyperlit.New(path)
	defer func() {
		if err := hl.Clear(); err != nil {
			fmt.Println(err)
		}
	}()

	cmd.InitCmds(hl)
	if err := cmd.Execute(ctx); err != nil {
		fmt.Println(err)
	}
}
