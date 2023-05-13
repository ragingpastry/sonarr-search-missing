package cmd

import (
	"github.com/ragingpastry/sonarr-search-missing/logger"
	"github.com/ragingpastry/sonarr-search-missing/search"
	"github.com/spf13/cobra"
)

var (
	debug    bool
	api_host string
	api_key  string
	log      = logger.NewLogger(false)
)

var rootCmd = &cobra.Command{
	Use:   "sonarr-search-missing [COMMAND]",
	Short: "Searches for missing episodes in Sonarr",
	Long:  "Uses the Sonarr API to look for missing episodes.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var log = logger.NewLogger(debug)
		search.Log = log
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug output")
	rootCmd.PersistentFlags().StringVarP(&api_host, "api-host", "", "", "API Endpoint to connect to.")
	rootCmd.PersistentFlags().StringVarP(&api_key, "api-key", "k", "", "API Key to use for authentication.")

	rootCmd.MarkPersistentFlagRequired("api-host")
	rootCmd.MarkPersistentFlagRequired("api-key")
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
