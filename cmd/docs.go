package cmd

import (
	"github.com/LeonidS635/HyperLit/internal/app/hyperlit"
	"github.com/spf13/cobra"
)

var (
	port int
	md   bool

	docs = &cobra.Command{
		Use:   "docs",
		Short: "Generate documentation",
	}
)

func docsCmd(hl *hyperlit.HyperLit) *cobra.Command {
	docs.RunE = func(cmd *cobra.Command, args []string) error {
		return hl.Docs(cmd.Context(), port, md)
	}
	return docs
}

func init() {
	docs.Flags().IntVarP(&port, "port", "p", 8123, "port to listen on")
	docs.Flags().BoolVarP(&md, "markdown", "m", true, "generate markdown documentation")
}
