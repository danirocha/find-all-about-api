package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
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
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Error getForecast -> http.Get. Using default value instead")

		body := fmt.Sprintf(`{
			"latitude": 52.52,
			"longitude": 13.419998,
			"generationtime_ms": 0.019073486328125,
			"utc_offset_seconds": 0,
			"timezone": "GMT",
			"timezone_abbreviation": "GMT",
			"elevation": 38,
			"daily_units": {
				"time": "iso8601",
				"temperature_2m_min": "°C",
				"temperature_2m_max": "°C"
			},
			"daily": {
				"time": [
					"%v"
				],
				"temperature_2m_min": [
					13.1
				],
				"temperature_2m_max": [
					26.3
				]
			}
		}`, time.Now().Format("2006-01-02"))
		resp = &http.Response{
			Status:        "200 OK",
			StatusCode:    200,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			ContentLength: int64(len(body)),
			// Request:       req,
			Header: make(http.Header, 0),
			Body:   io.NopCloser(bytes.NewBufferString(body)),
		}
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
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Error getImg -> http.Get. Using default value instead")

		body := `{
			"total_results": 10000,
			"page": 1,
			"per_page": 1,
			"photos": [
			  {
				"id": 3573351,
				"width": 3066,
				"height": 3968,
				"url": "https://www.pexels.com/photo/trees-during-day-3573351/",
				"photographer": "Lukas Rodriguez",
				"photographer_url": "https://www.pexels.com/@lukas-rodriguez-1845331",
				"photographer_id": 1845331,
				"avg_color": "#374824",
				"src": {
				  "original": "https://images.pexels.com/photos/3573351/pexels-photo-3573351.png",
				  "large2x": "https://images.pexels.com/photos/3573351/pexels-photo-3573351.png?auto=compress&cs=tinysrgb&dpr=2&h=650&w=940",
				  "large": "https://images.pexels.com/photos/3573351/pexels-photo-3573351.png?auto=compress&cs=tinysrgb&h=650&w=940",
				  "medium": "https://images.pexels.com/photos/3573351/pexels-photo-3573351.png?auto=compress&cs=tinysrgb&h=350",
				  "small": "https://images.pexels.com/photos/3573351/pexels-photo-3573351.png?auto=compress&cs=tinysrgb&h=130",
				  "portrait": "https://images.pexels.com/photos/3573351/pexels-photo-3573351.png?auto=compress&cs=tinysrgb&fit=crop&h=1200&w=800",
				  "landscape": "https://images.pexels.com/photos/3573351/pexels-photo-3573351.png?auto=compress&cs=tinysrgb&fit=crop&h=627&w=1200",
				  "tiny": "https://images.pexels.com/photos/3573351/pexels-photo-3573351.png?auto=compress&cs=tinysrgb&dpr=1&fit=crop&h=200&w=280"
				},
				"liked": false,
				"alt": "Brown Rocks During Golden Hour"
			  }
			],
			"next_page": "https://api.pexels.com/v1/search/?page=2&per_page=1&query=nature"
		  }
		  `
		resp = &http.Response{
			Status:        "200 OK",
			StatusCode:    200,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			ContentLength: int64(len(body)),
			// Request:       req,
			Header: make(http.Header, 0),
			Body:   io.NopCloser(bytes.NewBufferString(body)),
		}
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
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Error getIntro -> http.Get. Using default value instead")

		body := `{
			"pages": [
				{
					"pageid": 14072,
					"ns": 0,
					"index": 1,
					"type": "standard",
					"title": "History_of_Wikipedia",
					"displaytitle": "<span class=\"mw-page-title-main\">History of Wikipedia</span>",
					"namespace": {
						"id": 0,
						"text": ""
					},
					"wikibase_item": "Q6731",
					"titles": {
						"canonical": "History_of_Wikipedia",
						"normalized": "History of Wikipedia",
						"display": "<span class=\"mw-page-title-main\">History of Wikipedia</span>"
					},
					"thumbnail": {
						"source": "https://upload.wikimedia.org/wikipedia/commons/thumb/2/20/First_preserved_Main_Page_of_Wikipedia.jpeg/320px-First_preserved_Main_Page_of_Wikipedia.jpeg",
						"width": 320,
						"height": 391
					},
					"originalimage": {
						"source": "https://upload.wikimedia.org/wikipedia/commons/2/20/First_preserved_Main_Page_of_Wikipedia.jpeg",
						"width": 1280,
						"height": 1562
					},
					"lang": "en",
					"dir": "ltr",
					"revision": "1177457115",
					"tid": "0d2c05a0-5d5b-11ee-813c-4f1c66d7022d",
					"timestamp": "2023-09-27T17:26:44Z",
					"content_urls": {
						"desktop": {
							"page": "https://en.wikipedia.org/wiki/History_of_Wikipedia",
							"revisions": "https://en.wikipedia.org/wiki/History_of_Wikipedia?action=history",
							"edit": "https://en.wikipedia.org/wiki/History_of_Wikipedia?action=edit",
							"talk": "https://en.wikipedia.org/wiki/Talk:History_of_Wikipedia"
						},
						"mobile": {
							"page": "https://en.m.wikipedia.org/wiki/History_of_Wikipedia",
							"revisions": "https://en.m.wikipedia.org/wiki/Special:History/History_of_Wikipedia",
							"edit": "https://en.m.wikipedia.org/wiki/History_of_Wikipedia?action=edit",
							"talk": "https://en.m.wikipedia.org/wiki/Talk:History_of_Wikipedia"
						}
					},
					"extract": "Wikipedia, a free-content online encyclopedia written and maintained by a community of volunteers, began with its first edit on 15 January 2001, two days after the domain was registered. It grew out of Nupedia, a more structured free encyclopedia, as a way to allow easier and faster drafting of articles and translations.",
					"extract_html": "<p>Wikipedia, a free-content online encyclopedia written and maintained by a community of volunteers, began with its first edit on 15 January 2001, two days after the domain was registered. It grew out of Nupedia, a more structured free encyclopedia, as a way to allow easier and faster drafting of articles and translations.</p>",
					"normalizedtitle": "History of Wikipedia"
				}
			]
		}`
		resp = &http.Response{
			Status:        "200 OK",
			StatusCode:    200,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			ContentLength: int64(len(body)),
			// Request:       req,
			Header: make(http.Header, 0),
			Body:   io.NopCloser(bytes.NewBufferString(body)),
		}
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error getIntro -> io.ReadAll: %v", err)
	}
	var d IntroData
	if err := json.Unmarshal(data, &d); err != nil {
		log.Fatalf("Error getIntro -> json.Unmarshal: %v", err)
	}

	Coord = d.Pages[0].Coordinates

	return string(d.Pages[0].Extract)
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
