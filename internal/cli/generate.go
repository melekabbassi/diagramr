package cli

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/melekabbassi/diagramr/internal/config"
	"github.com/melekabbassi/diagramr/internal/output"
	base "github.com/melekabbassi/diagramr/internal/parser"
	golangparser "github.com/melekabbassi/diagramr/internal/parser/golang"
	"github.com/melekabbassi/diagramr/internal/renderer"
	"github.com/melekabbassi/diagramr/internal/studio"
	"github.com/spf13/cobra"
)

func newGenerateCmd() *cobra.Command {
	var (
		lang           string
		format         string
		outputPath     string
		includePrivate bool
		maxNodes       int
		noOpen         bool
	)

	cmd := &cobra.Command{
		Use:   "generate [path]",
		Short: "Generate diagram output from source code",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rootPath := "."
			if len(args) > 0 {
				rootPath = args[0]
			}

			cfg, err := config.Load()
			if err != nil {
				return err
			}

			if !cmd.Flags().Changed("lang") {
				lang = cfg.Language
			}
			if !cmd.Flags().Changed("format") {
				format = cfg.Format
			}
			if !cmd.Flags().Changed("max-nodes") {
				maxNodes = cfg.MaxNodes
			}

			if lang == "auto" {
				lang = detectLanguage(rootPath)
			}

			opts := base.Options{
				IncludePrivate: includePrivate,
			}

			var (
				content  string
				rendOpts renderer.Options
			)

			switch lang {
			case "go":
				g, parseErr := golangparser.New().Parse(rootPath, opts)
				if parseErr != nil {
					return fmt.Errorf("parse: %w", parseErr)
				}

				switch format {
				case "mermaid":
					rendOpts = renderer.Options{
						ShowMethods: true,
						ShowFields:  true,
						ShowPrivate: includePrivate,
						MaxNodes:    maxNodes,
					}
					content, err = renderer.NewClassRenderer().Render(g, rendOpts)
					if err != nil {
						return fmt.Errorf("render: %w", err)
					}
				case "json":
					content, err = output.RenderJSON(g)
					if err != nil {
						return fmt.Errorf("render json: %w", err)
					}
				default:
					return fmt.Errorf("unsupported format %q (available: mermaid, json)", format)
				}

			default:
				return fmt.Errorf("unsupported language %q (available: go)", lang)
			}

			// Write to file only when --output was explicitly given.
			if outputPath != "" {
				if err := output.Write(outputPath, content, cmd.OutOrStdout()); err != nil {
					return err
				}
			}

			if !noOpen && format == "mermaid" {
				// Studio displays and exports the diagram — skip stdout entirely.
				return studio.Serve(studio.Config{
					InitialMermaid: content,
					OutputPath:     outputPath,
					RenderConfig: studio.RenderConfig{
						RootPath:   rootPath,
						Lang:       lang,
						ParseOpts:  opts,
						RenderOpts: rendOpts,
					},
				})
			}

			// --no-open path: write to stdout when no --output file was given.
			if outputPath == "" {
				return output.Write("", content, cmd.OutOrStdout())
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&lang, "lang", "auto", "Source language (auto, go)")
	cmd.Flags().StringVar(&format, "format", "mermaid", "Output format (mermaid, json)")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path (default: stdout)")
	cmd.Flags().BoolVar(&includePrivate, "include-private", false, "Include unexported types and members")
	cmd.Flags().IntVar(&maxNodes, "max-nodes", 0, "Maximum number of nodes (0 = use config value)")
	cmd.Flags().BoolVar(&noOpen, "no-open", false, "Skip opening the browser after generation")

	return cmd
}

func detectLanguage(root string) string {
	var goCount int
	_ = filepath.WalkDir(root, func(_ string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if filepath.Ext(d.Name()) == ".go" {
			goCount++
		}
		return nil
	})
	if goCount > 0 {
		return "go"
	}
	return "go"
}
