package cmd

import (
	"context"

	"github.com/LeonidS635/HyperLit/internal/app/hyperlit"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hl",
	Short: "HyperLit CLI tool",
}

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}

func InitCmds(hl *hyperlit.HyperLit) {
	rootCmd.AddCommand(commitCmd(hl))
	rootCmd.AddCommand(docsCmd(hl))
}
