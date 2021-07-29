package cmd

import (
	"github.com/spf13/cobra"
)

var configFile string

var RootCmd = &cobra.Command{
	Use:   "win",
	Short: "Skill-based PvP arcade games powered by blockchain.",

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Run(APICmd, args)
	},
}

func init() {
	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.auth.yaml)")
}
