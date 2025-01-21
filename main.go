package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	//"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Points struct {
	Properties struct {
		GridID            string `json:"gridId"`
		GridX             int32  `json:"gridX"`
		GridY             int32  `json:"gridY"`
		ForecastZoneURL   string `json:"forecastZone"`
		HourlyForecastURL string `json:"forecastHourly"`
		ObservationURL    string `json:"observationStations"`
		RelativeLocation  struct {
			Properties struct {
				City  string `json:"city"`
				State string `json:"state"`
			} `json: "properties"`
		} `json: "relativeLocation"`
	} `json: "properties"`
}

type Stations struct {
	Features []struct {
		Properties struct {
			StationID   string `json:"stationIdentifier"`
			StationName string `json:"name"`
		} `json:"properties"`
	} `json:"features"`
}

type Observations struct {
	Properties struct {
		Description string `json:"textDescription"`
		Temperature struct {
			Value float64 `json:"value"`
		} `json:"temperature"`
		Humidity struct {
			Value float64 `json:"value"`
		} `json:"relativeHumidity"`
	} `json:"properties"`
}

type HourlyForecast struct {
	Properties struct {
		Periods []struct {
			Number        int32  `json:"number"`
			StartTime     string `json:"startTime"`
			EndTime       string `json:"endTime"`
			Temperature   int    `json:"temperature"`
			Precipitation struct {
				Unit  string `json:"unitCode"`
				Value int32  `json:"value"`
			} `json:"propabilityOfPrecipitation"`
			Humidity struct {
				Unit  string `json:"unitCode"`
				Value int32  `json:"value"`
			} `json:"relativeHumidity"`
			ShortForecast string `json:"shortForecast"`
		} `json:"periods"`
	} `json:"properties"`
}

type Alerts struct {
	Features []struct {
		Properties struct {
			Severity      string `json:"severity"`
			EffectiveTime string `json:"effective"`
			ExpireTime    string `json:"expires"`
			Event         string `json:"event"`
			Sender        string `json:"senderName"`
			Headline      string `json:"headline"`
			Description   string `json:"description"`
			Instruction   string `json:"instruction"`
		} `json:"properties"`
	} `json:"features"`
}

func main() {
	/*
	  STAGE 1: CONVERT LATITUDE & LONGITUDE TO NWS GRID VALUES

	  This stage will:
	    - prompt user for location (zip code/known location)
	    - get the NWS grid points for the corresponding lat/long tuple
	    - parse the data and store the important values in vars, which are used to build the data
	      we will print on the console for users

	*/

	// TODO : prompt for zip code and convert to lat/long
	lat := 29.7725
	long := -95.6201

	//if len(os.Args) >= 3 {
	//	lat = os.Args[1]
	//	long = os.Args[2]
	//}

	res, err := http.Get(fmt.Sprintf("https://api.weather.gov/points/%f,%f", lat, long))
	// check request errors
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	// check response errors
	if res.StatusCode != 200 {
		panic("Weather API not available")
	}

	// attempt to parse the response data
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// unmarshall the data into our custom struct
	var points Points
	err = json.Unmarshal(body, &points)
	if err != nil {
		panic(err)
	}

	// storing values we need into vars
	currentCity := points.Properties.RelativeLocation.Properties.City
	currentState := points.Properties.RelativeLocation.Properties.State
	// the zone ID is also useful, but it's stuck in a URL.  Let's parse it.
	u, err := url.Parse(points.Properties.ForecastZoneURL)
	if err != nil {
		panic(err)
	}
	forecastZonePath := strings.Split(u.Path, "/")
	forecastZone := forecastZonePath[len(forecastZonePath)-1]

	/*
	  STAGE 2: GET CURRENT WEATHER INFORMATION

	  This stage will:
	    - use the NWS grid data to locate a nearby station
	    - query the station for its latest weather data ("observations")
	    - parse the data and store the important values in vars, which are used to build the data
	      we will print on the console for users

	*/
	// the original grid API gives us a station URL, let's use that rather than build our own
	stationObservationURL := points.Properties.ObservationURL
	// make the request
	res, err = http.Get(stationObservationURL)
	// check for request errors
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	// check for response errors
	if res.StatusCode != 200 {
		panic("Weather API not available")
	}

	// attempt to parse response data
	body, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// unmarshall data to our struct
	var stations Stations
	err = json.Unmarshal(body, &stations)
	if err != nil {
		panic(err)
	}

	// store the FIRST station ID and Name for later...
	stationID := stations.Features[0].Properties.StationID
	stationName := stations.Features[0].Properties.StationName

	res, err = http.Get(fmt.Sprintf("https://api.weather.gov/stations/%s/observations/latest", stationID))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Station API not available")
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	// unmarshall data to our struct
	var observations Observations
	err = json.Unmarshal(body, &observations)
	if err != nil {
		panic(err)
	}

	// the station API returns degC, not degF, so let's convert:
	currentTempC := observations.Properties.Temperature.Value
	currentTempF := (currentTempC * (9 / 5)) + 32
	currentHumidity := observations.Properties.Humidity.Value
	currentCondition := observations.Properties.Description

	/*
	  STAGE 3: GET HOURLY FORECAST DATA

	  This stage will:
	    - use the NWS grid data to query for houry forecasts
	    - parse the data and store the important values in vars, which are used to build the data
	      we will print on the console for users

	*/

	// the original grid API gives us a forecast URL, let's use that rather than build our own
	hourlyForecastURL := points.Properties.HourlyForecastURL

	// make the request
	res, err = http.Get(hourlyForecastURL)
	// check for request errors
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	// check for response errors
	if res.StatusCode != 200 {
		panic("Weather API not available")
	}

	// attempt to parse response data
	body, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// unmarshall data to our struct
	var hourlyForecast HourlyForecast
	err = json.Unmarshal(body, &hourlyForecast)
	if err != nil {
		panic(err)
	}

	/*
	  STAGE 4: GET CURRENT ALERTS

	  This stage will:
	    - use the NWS grid data to query for active alerts for the current location
	    - parse the data and store the important values in vars, which are used to build the data
	      we will print on the console for users

	*/
	res, err = http.Get(fmt.Sprintf("https://api.weather.gov/alerts/active/zone/%s", forecastZone))
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		panic("Alert API not available")
	}

	// attempt to parse response data
	body, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// unmarshall data to our struct
	var alerts Alerts
	err = json.Unmarshal(body, &alerts)
	if err != nil {
		panic(err)
	}

	/*
			  STAGE 5: PRINT WEATHER INFORMATION

			  This stage will:
			    - define custom print styles (courtesy of fatih/color)
		      - build message contents with forecast data
		      - (TODO) add formatting based on any thresholds
		      - print the weather data to the console

	*/

	// Define custom styles for printing
	boldWhite := color.RGB(199, 208, 244).Add(color.Bold)
	italicWhite := color.RGB(199, 208, 244).Add(color.Italic)
	italicBlack := color.New(color.FgBlack).Add(color.Italic)
	boldRed := color.New(color.FgRed).Add(color.Bold)

	// MAKE THIS BOLD
	boldWhite.Println(fmt.Sprintf(
		"%s, %s\n",
		currentCity,
		currentState,
	))

	// ALERTS (SHOULD BE RED & BOLD)
	boldRed.Println("CURRENT ALERTS")
	// instructions should be italic & unbolded
	for _, feature := range alerts.Features {
		startTime, err := time.Parse(time.RFC3339, feature.Properties.EffectiveTime)
		if err != nil {
			panic(err)
		}
		endTime, err := time.Parse(time.RFC3339, feature.Properties.ExpireTime)
		if err != nil {
			panic(err)
		}
		if endTime.Before(time.Now()) {
			continue
		} else {
			message := fmt.Sprintf("%s (%s)\nStart: %s\tEnd: %s\n\n",
				feature.Properties.Event,
				feature.Properties.Sender,
				startTime.Format("02 Jan 2006 15:04:05"),
				endTime.Format("02 Jan 2006 15:04:05"),
			)
			color.Red(message)
		}
	}

	// CURRENT WEATHER
	boldWhite.Println("CURRENT CONDITIONS")
	// values
	italicWhite.Println("Temperature\tHumidity\tSummary")
	fmt.Printf("%.0f°F\t\t%.1f%%\t\t%s\n",
		currentTempF,
		currentHumidity,
		currentCondition,
	)
	// print the source station info for extra data
	italicBlack.Printf("%s (%s)\n\n",
		stationID,
		stationName,
	)

	// FORECAST
	boldWhite.Println("FORECAST")
	italicWhite.Println("Time\t\tTemperature\tPrecip %\tHumidity %\tDescription")
	for _, period := range hourlyForecast.Properties.Periods {
		currentDay := time.Now().Day()
		forecastDateTime, err := time.Parse(time.RFC3339, period.StartTime)
		if err != nil {
			panic(err)
		}
		forecastDay := forecastDateTime.Day()
		forecastTemp := period.Temperature
		forecastPrecip := period.Precipitation.Value
		forecastHumid := period.Humidity.Value
		forecastDesc := period.ShortForecast
		// exclude anything except the current day
		if forecastDay > currentDay {
			continue
		} else {
			fmt.Printf("%s\t\t%d°F\t\t%d%%\t\t%d%%\t\t%s\n",
				forecastDateTime.Format("15:04"),
				forecastTemp,
				forecastPrecip,
				forecastHumid,
				forecastDesc,
			)
		}
	}
}
