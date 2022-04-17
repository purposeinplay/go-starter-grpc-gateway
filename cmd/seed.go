package cmd

import (
	"fmt"

	"github.com/purposeinplay/go-starter-grpc-gateway/internal/adapters/psql"
	"github.com/purposeinplay/go-starter-grpc-gateway/internal/common/config"
	"github.com/spf13/cobra"
)

// CMD executes a Seed Command.
var CMD = &cobra.Command{
	Use: "seed",
	Long: "Seed database strucutures. " +
		"This will create new tables and add missing columns and indexes.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig(cmd)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		_, err = psql.Connect(cfg)
		if err != nil {
			return fmt.Errorf("connect db: %w", err)
		}

		return nil
	},
}
