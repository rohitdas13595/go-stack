package cli

import (
	"os"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
)

var servePort int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run dev server (go run ./cmd/server)",
	Run: func(cmd *cobra.Command, args []string) {
		_ = os.Setenv("PORT", strconv.Itoa(servePort))
		c := exec.Command("go", "run", "./cmd/server")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Env = os.Environ()
		exitErr(c.Run())
	},
}

func init() {
	serveCmd.Flags().IntVar(&servePort, "port", 3000, "listen port")
}
