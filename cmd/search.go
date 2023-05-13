package cmd

import (
	"fmt"

	"github.com/ragingpastry/sonarr-search-missing/search"
	"github.com/spf13/cobra"
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

var (
	searchAll bool
	series    string
	seasons   []string
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for missing episodes",
	Long:  "Search for missing episodes",
	PreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Lookup("seasons").Changed {
			cmd.MarkFlagRequired("series")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if searchAll {
			log.Info("Searching all series")
			api_key, _ := cmd.Flags().GetString("api-key")
			api_host, _ := cmd.Flags().GetString("api-host")
			c := starr.New(api_key, api_host, 0)
			s := sonarr.New(c)
			search.SearchAll(s)
		} else if series != "" {
			log.Info(fmt.Sprintf("Searching series %s", series))
			api_key, _ := cmd.Flags().GetString("api-key")
			api_host, _ := cmd.Flags().GetString("api-host")
			c := starr.New(api_key, api_host, 0)
			s := sonarr.New(c)
			series := search.GetSeries(s, series)
			if series == nil {
				log.Error(fmt.Sprintf("Series %s not found", series.Title))
			}
			if len(seasons) > 0 {
				search.SearchSeason(s, series, seasons)
			} else {
				search.SearchSeries(s, series)
			}
		} else {
			cmd.Help()
		}

	},
}

func bindSearchFlags() {
	searchFlags := searchCmd.Flags()

	searchFlags.BoolVarP(&searchAll, "all", "a", false, "Search all series")
	searchFlags.StringVarP(&series, "series", "s", "", "Search a specific series")
	searchFlags.StringSliceVarP(&seasons, "seasons", "", []string{}, "Search specific seasons in a series")
}

func init() {
	rootCmd.AddCommand(searchCmd)
	bindSearchFlags()
}
