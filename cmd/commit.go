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
		RunE: func(cmd *cobra.Command, args []string) error {
			if needToSave, err := hl.CommitFirstStep(cmd.Context()); needToSave && err == nil {
				var ans string
				fmt.Print("Do you want to save changes? [y/n] ")
				_, err = fmt.Scan(&ans)
				if err != nil {
					return err
				}
				ans = strings.ToLower(ans)

				if ans == "y" || ans == "yes" {
					return hl.CommitSecondStep(cmd.Context())
				}
				return nil
			} else {
				return err
			}
		},
	}
}
