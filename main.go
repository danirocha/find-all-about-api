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

type Page struct {
	Extract string
}

type Data struct {
	Pages []Page
}

type Location struct {
	Name  string
	Intro string
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
	var d Data
	if err := json.Unmarshal(data, &d); err != nil {
		log.Fatal(err)
	}

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

	s, err := json.Marshal(Location{location, intro})

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
