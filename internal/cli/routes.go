package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var routesCmd = &cobra.Command{
	Use:   "routes",
	Short: "List routes (use app.Routes() in code or wire debug endpoint)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Tip: call app.Routes() from your bootstrap or add a debug handler that prints JSON.")
	},
}
