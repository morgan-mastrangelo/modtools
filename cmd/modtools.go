package cmd

import (
	"os"

	"github.com/muesli/coral"
)

var rootCmd = &coral.Command{
	Use: "modtools",
	PersistentPreRun: func(cmd *coral.Command, args []string) {
		cmd.SilenceUsage = true
	},
}

func Execute() {
	rootCmd.PersistentFlags().Bool("direct-only", false, "work only on direct dependencies")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func getDirectOnly(cmd *coral.Command) bool {
	return cmd.Flags().Lookup("direct-only").Value.String() != "false"
}
