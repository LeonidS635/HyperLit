package cmd

import (
	"github.com/LeonidS635/HyperLit/internal/app/hyperlit"
	"github.com/spf13/cobra"
)

func docsCmd(hl *hyperlit.HyperLit) *cobra.Command {
	return &cobra.Command{
		Use:   "docs",
		Short: "Generate documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			return hl.Docs(cmd.Context())
		},
	}
}
