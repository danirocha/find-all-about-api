package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func FindAllAbout(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatalf("An error occurred: %v", err)
		return
	}

	location := r.FormValue("location")
	fmt.Println("SEARCH ALL ABOUT ", location)

	io.WriteString(w, location)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./welcome")))
	http.HandleFunc("/find-all-about", FindAllAbout)

	if err := http.ListenAndServe(":3031", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
		os.Exit(1)
	}
}
