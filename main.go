package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetHello(w http.ResponseWriter, r *http.Request) {
	s := "Welcome to POI api!"

	io.WriteString(w, s)
}

func main() {
	http.HandleFunc("/", GetHello)
	err := http.ListenAndServe(":3031", nil)

	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}
