package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
)

func getStars(w http.ResponseWriter, r *http.Request) {
	client := github.NewClient(nil)

	// list all organizations for user "willnorris"
	orgs, _, err := client.Organizations.List(context.Background(), "willnorris", nil)

	if err != nil {
		fmt.Printf("Cannot call go-github %v", err)
	}
	for i, organization := range orgs {
		fmt.Printf("%v. %v\n", i+1, organization.GetLogin())
	}

	fmt.Fprintf(w, "Done")
}

func main() {
	fmt.Println("... start server")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/getStars", getStars)
	log.Fatal(http.ListenAndServe(":8080", router))
}
