package list

import (
	"fmt"

	"github.com/ragingpastry/sonarr-search-missing/logger"
	"golift.io/starr/sonarr"
)

var Log *logger.Logger

// / List monitored series
// / @param sonarrConnection: Sonarr connection
// / @return: void
func ListMonitored(sonarrConnection *sonarr.Sonarr) {
	allSeries, err := sonarrConnection.GetAllSeries()
	if err != nil {
		Log.Error("Error getting all series")
	}
	for _, series := range allSeries {
		if series.Monitored {
			Log.Info(fmt.Sprintf("Monitored series: %s", series.Title))
		}
	}
}

func ListAllMissing(sonarrConnection *sonarr.Sonarr) int {
	allSeries, err := sonarrConnection.GetAllSeries()
	if err != nil {
		Log.Error("Error getting all series")
	}
	totalMissing := 0
	for _, series := range allSeries {
		if series.Statistics.PercentOfEpisodes < 100 && series.Monitored {
			totalMissing = (series.Statistics.TotalEpisodeCount - series.Statistics.EpisodeCount) + totalMissing
		}
	}
	return totalMissing
}

func ListMissing(sonarrConnection *sonarr.Sonarr, seriesName string) (int, error) {
	allSeries, err := sonarrConnection.GetAllSeries()
	for series := range allSeries {
		if allSeries[series].Title == seriesName {
			if allSeries[series].Statistics.PercentOfEpisodes < 100 && allSeries[series].Monitored {
				return allSeries[series].Statistics.TotalEpisodeCount - allSeries[series].Statistics.EpisodeCount, nil
			}
		}
	}
	return 0, err
}
