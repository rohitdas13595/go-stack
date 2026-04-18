package cli

import (
	"fmt"
	"os"

	"github.com/rohitdas13595/go-stack/config"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Print resolved config from config/*.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		env := os.Getenv("APP_ENV")
		if env == "" {
			env = "development"
		}
		_, err := config.Load("config/app.yaml", "config/"+env+".yaml")
		if err != nil {
			exitErr(err)
			return
		}
		if g := config.Global(); g != nil {
			fmt.Printf("app.name=%q\n", g.String("app.name"))
			fmt.Printf("server.port=%d\n", g.Int("server.port"))
		}
	},
}
