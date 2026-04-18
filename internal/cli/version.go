package cli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print diagramr version information",
		Run: func(cmd *cobra.Command, args []string) {
			_, _ = fmt.Fprintf(
				cmd.OutOrStdout(),
				"diagramr version %s\ncommit: %s\nbuild_date: %s\ngo_version: %s\n",
				version,
				commit,
				date,
				runtime.Version(),
			)
		},
	}
}
