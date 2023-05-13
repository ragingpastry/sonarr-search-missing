package search

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ragingpastry/sonarr-search-missing/logger"
	"golift.io/starr/sonarr"
)

var Log *logger.Logger

// / Search for missing episodes
// / @param sonarrConnection: Sonarr connection
// / @return: void
// / @throws: error
func SearchAll(sonarrConnection *sonarr.Sonarr) {
	allSeries, err := sonarrConnection.GetAllSeries()
	if err != nil {
		Log.Error("Error getting all series")
	}
	for _, series := range allSeries {
		if series.Statistics.PercentOfEpisodes < 100 && series.Monitored {
			if series.Title == "The Witcher" {
				Log.Info(fmt.Sprintf("Search for missing episodes for %s", series.Title))
				s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
				s.Suffix = fmt.Sprintf(" Searching for missing episodes for %s", series.Title)
				commandRequest := &sonarr.CommandRequest{
					Name:     "MissingEpisodeSearch",
					SeriesID: series.ID,
				}
				response, err := sonarrConnection.SendCommand(commandRequest)
				if err != nil {
					Log.Error(fmt.Sprintf("Error sending command: %s", err))
				}
				s.Start()
				waitForCommandToFinish(response, sonarrConnection)
				s.Stop()
				Log.Info(fmt.Sprintf("Search for missing episodes for %s complete", series.Title))
			}

		}
	}
}

func SearchSeries(sonarrConnection *sonarr.Sonarr, series *sonarr.Series) {
	if series.Statistics.PercentOfEpisodes < 100 && series.Monitored {
		Log.Info(fmt.Sprintf("Search for missing episodes for %s", series.Title))
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Suffix = fmt.Sprintf(" Searching for missing episodes for %s", series.Title)
		commandRequest := &sonarr.CommandRequest{
			Name:     "MissingEpisodeSearch",
			SeriesID: series.ID,
		}
		response, err := sonarrConnection.SendCommand(commandRequest)
		if err != nil {
			Log.Error(fmt.Sprintf("Error sending command: %s", err))
		}
		s.Start()
		waitForCommandToFinish(response, sonarrConnection)
		s.Stop()
		Log.Info(fmt.Sprintf("Search for missing episodes for %s complete", series.Title))
	}
}

func SearchSeason(sonarrConnection *sonarr.Sonarr, series *sonarr.Series, seasons []string) {
	for _, season := range seasons {
		Log.Info(fmt.Sprintf("Search for missing episodes for %s : Season %s", series.Title, season))
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Suffix = fmt.Sprintf(" Searching for missing episodes for %s : Season %s", series.Title, season)

		commandRequest := &sonarr.CommandRequest{
			Name:     "MissingEpisodeSearch",
			SeriesID: series.ID,
		}
		response, err := sonarrConnection.SendCommand(commandRequest)
		if err != nil {
			Log.Error(fmt.Sprintf("Error sending command: %s", err))
		}
		s.Start()
		waitForCommandToFinish(response, sonarrConnection)
		s.Stop()
		Log.Info(fmt.Sprintf("Search for missing episodes for %s : Season %s complete", series.Title, season))
	}
}

func GetSeries(sonarrConnection *sonarr.Sonarr, seriesName string) *sonarr.Series {
	allSeries, err := sonarrConnection.GetAllSeries()
	if err != nil {
		Log.Error("Error getting all series")
	}
	for _, series := range allSeries {
		if series.Title == seriesName {
			return series
		}
	}
	return nil
}

// / Wait for command to finish
// / @param commandResponse: Command response
// / @param sonarrConnection: Sonarr connection
// / @return: void
func waitForCommandToFinish(commandResponse *sonarr.CommandResponse, sonarrConnection *sonarr.Sonarr) {

	commandStatusChan := make(chan *sonarr.CommandResponse)

	go func() {
		for {
			commandResponse, err := sonarrConnection.GetCommandStatus(commandResponse.ID)
			if err != nil {
				panic(err)
			}
			commandStatusChan <- commandResponse

			if commandResponse.Status == "completed" {
				return
			}

			time.Sleep(5 * time.Second)

		}
	}()

	for {
		select {
		case commandResponse := <-commandStatusChan:
			Log.Debug(fmt.Sprintf("Command status: %s\n", commandResponse.Status))
			if commandResponse.Status == "completed" {
				return
			}
		}
	}

}
