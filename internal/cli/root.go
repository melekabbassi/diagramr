package cli

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diagramr",
		Short: "Generate Mermaid diagrams from source code",
	}
	cmd.AddCommand(newGenerateCmd(), newInitCmd(), newValidateCmd(), newVersionCmd())
	return cmd
}
