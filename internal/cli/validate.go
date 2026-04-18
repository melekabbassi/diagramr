package cli

import (
	"fmt"

	"github.com/melekabbassi/diagramr/internal/config"
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration and print a summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			if cfg.MaxNodes <= 0 {
				return fmt.Errorf("invalid max_nodes: must be > 0")
			}

			_, _ = fmt.Fprintf(
				cmd.OutOrStdout(),
				"Config valid\nlanguage: %s\nformat: %s\nmax_nodes: %d\n",
				cfg.Language,
				cfg.Format,
				cfg.MaxNodes,
			)

			return nil
		},
	}
}
