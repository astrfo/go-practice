package main

import (
	"fmt"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to My HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func main() {
	fmt.Println("AccessURL: http://localhost:8080/")
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":8080", nil)
}
