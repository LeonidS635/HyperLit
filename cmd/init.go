package cmd

import (
	"github.com/LeonidS635/HyperLit/internal/app/hyperlit"
	"github.com/spf13/cobra"
)

func initCmd(hl *hyperlit.HyperLit) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "some info",
		Run: func(cmd *cobra.Command, args []string) {
			hl.Init(cmd.Context())
		},
	}
}
