package cli

import (
	"encoding/json"
	"os"

	"github.com/rohitdas13595/go-stack/openapi"
	"github.com/spf13/cobra"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Emit stub OpenAPI JSON to stdout",
	Run: func(cmd *cobra.Command, args []string) {
		spec := openapi.FromRoutes([][2]string{
			{"GET", "/health"},
		})
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(spec); err != nil {
			exitErr(err)
		}
	},
}
