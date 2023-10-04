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

type Img struct {
	Url string
}

type ImgData struct {
	Photos []Img
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
	Img      string
}

var Coord Coordinates

func getForecast(lat, lon float64) []Forecast {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?daily=temperature_2m_max,temperature_2m_min&forecast_days=1&latitude=%v&longitude=%v", lat, lon)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error getForecast -> http.Get: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error getForecast -> io.ReadAll: %v", err)
	}

	var d ForecastData
	if err := json.Unmarshal(data, &d); err != nil {
		log.Fatalf("Error getForecast -> json.Unmarshal: %v", err)
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

const SERVICE_PEXELS_API_KEY = ""

func getImg(location string) string {
	url := fmt.Sprintf("https://api.pexels.com/v1/search?query=%v&per_page=1", location)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error getImg -> http.Get: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error getImg -> io.ReadAll: %v", err)
	}

	var d ImgData
	if err := json.Unmarshal(data, &d); err != nil {
		log.Fatalf("Error getImg -> json.Unmarshal: %v", err)
	}

	return d.Photos[0].Url
}

func getIntro(location string) string {
	url := "https://en.wikipedia.org/api/rest_v1/page/related/"
	re := regexp.MustCompile(`\s+`)
	url += re.ReplaceAllString(location, "_")

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error getIntro -> http.Get: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error getIntro -> io.ReadAll: %v", err)
	}
	if string(data) != "" {
		var d IntroData
		if err := json.Unmarshal(data, &d); err != nil {
			log.Fatalf("Error getIntro -> json.Unmarshal: %v", err)
		}

		Coord = d.Pages[0].Coordinates

		return string(d.Pages[0].Extract)
	} else {
		return ""
	}
}

func FindAllAbout(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatalf("Error r.ParseForm: %v", err)
		return
	}

	location := r.FormValue("location")
	fmt.Println("SEARCH ALL ABOUT ", location)
	intro := getIntro(location)
	forecast := getForecast(Coord.Lat, Coord.Lon)
	img := getImg(location)

	s, err := json.Marshal(Location{location, intro, forecast, img})

	if err != nil {
		log.Fatalf("Error json.Marshal: %v", err)
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
