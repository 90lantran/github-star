package main

import (
	"fmt"

	_ "github.com/pdrum/swagger-automation/docs" // This line is necessary for go-swagger to find your docs!
)

func main() {
	fmt.Println("... starts server")
	// router := mux.NewRouter().StrictSlash(true)
	// router.HandleFunc("/get-stars", getStars).Methods("POST")
	// //router.HandleFunc("/getStars", getStars)
	// log.Fatal(http.ListenAndServe(":8080", router))
	server := server{}
	server.startServer()
}
