package cli

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/rohitdas13595/go-stack/db"
	"github.com/rohitdas13595/go-stack/migrate"
"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database commands",
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run pending SQL migrations",
	Run: func(cmd *cobra.Command, args []string) {
		exitErr(runMigrate(false))
	},
}

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback last migration batch",
	Run: func(cmd *cobra.Command, args []string) {
		exitErr(runMigrate(true))
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show applied migrations",
	Run: func(cmd *cobra.Command, args []string) {
		exitErr(showStatus())
	},
}

func init() {
	dbCmd.AddCommand(migrateCmd, rollbackCmd, statusCmd)
}

func openAppDB() (*sql.DB, string, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "file:./storage/app.db"
	}
	sdb, err := db.OpenSQLite(dsn)
	if err != nil {
		return nil, "", err
	}
	return sdb, migrate.Dialect("sqlite"), nil
}

func runMigrate(rollback bool) error {
	ctx := context.Background()
	sqldb, dialect, err := openAppDB()
	if err != nil {
		return err
	}
	defer sqldb.Close()
	ms, err := migrate.LoadDir("db/migrations")
	if err != nil {
		return err
	}
	if rollback {
		return migrate.Rollback(ctx, sqldb, dialect, ms)
	}
	return migrate.Migrate(ctx, sqldb, dialect, ms)
}

func showStatus() error {
	ctx := context.Background()
	sqldb, _, err := openAppDB()
	if err != nil {
		return err
	}
	defer sqldb.Close()
	vs, err := migrate.Status(ctx, sqldb)
	if err != nil {
		return err
	}
	for _, v := range vs {
		fmt.Println(v)
	}
	return nil
}
