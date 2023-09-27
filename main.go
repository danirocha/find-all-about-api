package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
)

type ScaleUnit struct {
	Unit string `json:"temperature_2m_max"`
}

type Daily struct {
	MinTempList  []float64 `json:"temperature_2m_min"`
	MaxTempList  []float64 `json:"temperature_2m_max"`
	DateTempList []string  `json:"time"`
}

type ForecastData struct {
	ScaleUnit ScaleUnit `json:"daily_units"`
	Daily     Daily
}

type Coordinates struct {
	Lat float64
	Lon float64
}

type Page struct {
	Extract     string
	Coordinates Coordinates
}

type IntroData struct {
	Pages []Page
}

type Forecast struct {
	Date           string
	MinTemperature float64
	MaxTemperature float64
	ScaleUnit      string
}

type Location struct {
	Name     string
	Intro    string
	Forecast []Forecast
}

var Coord Coordinates

func getForecast(lat, lon float64) []Forecast {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?daily=temperature_2m_max,temperature_2m_min&forecast_days=1&latitude=%v&longitude=%v", lat, lon)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var d ForecastData
	if err := json.Unmarshal(data, &d); err != nil {
		log.Fatal(err)
	}

	return []Forecast{
		{
			d.Daily.DateTempList[0],
			d.Daily.MinTempList[0],
			d.Daily.MaxTempList[0],
			d.ScaleUnit.Unit,
		},
	}
}

func getIntro(location string) string {
	url := "https://en.wikipedia.org/api/rest_v1/page/related/"
	re := regexp.MustCompile(`\s+`)
	url += re.ReplaceAllString(location, "_")

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var d IntroData
	if err := json.Unmarshal(data, &d); err != nil {
		log.Fatal(err)
	}

	Coord = d.Pages[0].Coordinates

	return string(d.Pages[0].Extract)
}

func FindAllAbout(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
		return
	}

	location := r.FormValue("location")
	fmt.Println("SEARCH ALL ABOUT ", location)
	intro := getIntro(location)
	forecast := getForecast(Coord.Lat, Coord.Lon)

	s, err := json.Marshal(Location{location, intro, forecast})

	if err != nil {
		log.Fatal(err)
	}

	io.WriteString(w, string(s))
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./welcome")))
	http.HandleFunc("/find-all-about", FindAllAbout)

	if err := http.ListenAndServe(":3031", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
		os.Exit(1)
	}
}
