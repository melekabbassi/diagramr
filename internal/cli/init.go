package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/melekabbassi/diagramr/internal/config"
	"github.com/spf13/cobra"
)

const defaultConfigFile = "diagramr.config.yaml"

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Create a default diagramr config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(defaultConfigFile); err == nil {
				return fmt.Errorf("config file already exists: %s", defaultConfigFile)
			} else if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("check config file: %w", err)
			}

			cfg := config.Default()
			content := fmt.Sprintf("language: %s\nformat: %s\nmax_nodes: %d\n", cfg.Language, cfg.Format, cfg.MaxNodes)

			if err := os.WriteFile(defaultConfigFile, []byte(content), 0o644); err != nil {
				return fmt.Errorf("write config file: %w", err)
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", defaultConfigFile)
			return nil
		},
	}
}
