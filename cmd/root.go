package cmd

import (
	"fmt"
	"os"

	"github.com/LeonidS635/HyperLit/internal/app/hyperlit"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hl",
	Short: "HyperLit CLI tool",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func InitCmds(hl *hyperlit.HyperLit) {
	rootCmd.AddCommand(initCmd(hl))
	rootCmd.AddCommand(commitCmd(hl))
	rootCmd.AddCommand(docsCmd(hl))
}
