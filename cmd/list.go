package cmd

import (
	"fmt"

	"github.com/ragingpastry/sonarr-search-missing/list"
	"github.com/spf13/cobra"
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List missing episodes",
	Long:  "List missing episodes",
}

var ListMissingCmd = &cobra.Command{
	Use:   "missing",
	Short: "List missing episodes",
	Long:  "List missing episodes",
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired("series")
	},
	Run: func(cmd *cobra.Command, args []string) {
		api_key, _ := cmd.Flags().GetString("api-key")
		api_host, _ := cmd.Flags().GetString("api-host")
		series, _ := cmd.Flags().GetString("series")
		c := starr.New(api_key, api_host, 0)
		s := sonarr.New(c)
		if series != "" {
			log.Info(fmt.Sprintf("Listing missing episodes for series %s", series))
			numMissing, err := list.ListMissing(s, series)
			if err != nil {
				log.Error(fmt.Sprintf("Error listing missing episodes: %s", err))
			}
			log.Info(fmt.Sprintf("Found %d missing episodes for series %s", numMissing, series))
		} else {
			numMissing := list.ListAllMissing(s)

			log.Info(fmt.Sprintf("Found %d missing episodes", numMissing))
		}
	},
}

var ListMonitoredCmd = &cobra.Command{
	Use:   "monitored",
	Short: "List monitored series",
	Long:  "List monitored series",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Listing monitored series")
		api_key, _ := cmd.Flags().GetString("api-key")
		api_host, _ := cmd.Flags().GetString("api-host")
		c := starr.New(api_key, api_host, 0)
		s := sonarr.New(c)
		list.ListMonitored(s)

	},
}

func bindListFlags() {
	listFlags := listCmd.Flags()
	listMissingFlags := ListMissingCmd.Flags()

	listFlags.BoolVarP(&searchAll, "all", "a", false, "List all series")
	listMissingFlags.StringVarP(&series, "series", "s", "", "Series to search")
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(ListMonitoredCmd)
	listCmd.AddCommand(ListMissingCmd)
	bindListFlags()
}
