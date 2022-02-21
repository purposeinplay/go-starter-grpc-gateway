package cmd

import (
	"fmt"

	"github.com/purposeinplay/go-starter-grpc-gateway/cmd/migrate"
	"github.com/purposeinplay/go-starter-grpc-gateway/cmd/seed"
	"github.com/purposeinplay/go-starter-grpc-gateway/cmd/server"
	"github.com/spf13/cobra"
)

var configFile string

func init() {
	RootCmd.AddCommand(
		server.CMD,
		seed.CMD,
		migrate.CMD,
	)

	RootCmd.PersistentFlags().StringVar(
		&configFile,
		"config",
		"",
		"config file (default is $HOME/.auth.yaml)",
	)
}

// RootCmd represents the application root command.
var RootCmd = &cobra.Command{
	Use:   "win",
	Short: "Skill-based PvP arcade games powered by blockchain.",

	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.RunE(server.CMD, args)
		if err != nil {
			return fmt.Errorf("run server cmd: %w", err)
		}

		return nil
	},
}
