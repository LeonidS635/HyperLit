package cmd

import (
	"fmt"

	"github.com/LeonidS635/HyperLit/internal/app/hyperlit"
	"github.com/spf13/cobra"
)

func commitCmd(hl *hyperlit.HyperLit) *cobra.Command {
	return &cobra.Command{
		Use:   "commit",
		Short: "Commit changes",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Context())
			hl.Commit(cmd.Context())
		},
	}
}
