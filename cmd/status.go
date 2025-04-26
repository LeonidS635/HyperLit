package cmd

import (
	"github.com/LeonidS635/HyperLit/internal/app/hyperlit"
	"github.com/spf13/cobra"
)

func statusCmd(hl *hyperlit.HyperLit) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "show files status",
		Run: func(cmd *cobra.Command, args []string) {
			hl.Status(cmd.Context())
		},
	}
}
