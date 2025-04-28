package cmd

import (
	"fmt"
	"strings"

	"github.com/LeonidS635/HyperLit/internal/app/hyperlit"
	"github.com/spf13/cobra"
)

func commitCmd(hl *hyperlit.HyperLit) *cobra.Command {
	return &cobra.Command{
		Use:   "commit",
		Short: "Commit changes",
		Run: func(cmd *cobra.Command, args []string) {
			hl.CommitFirstStep(cmd.Context())

			var ans string
			fmt.Print("Do you want to save changes? [y/n] ")
			fmt.Scan(&ans)
			ans = strings.ToLower(ans)

			if ans == "y" || ans == "yes" {
				hl.CommitSecondStep(cmd.Context())
			}
		},
	}
}
