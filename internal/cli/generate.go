package cli

import "github.com/spf13/cobra"

func newGenerateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate [path]",
		Short: "Generate diagram output",
		RunE:  func(cmd *cobra.Command, args []string) error { return nil },
	}
}
