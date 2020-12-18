package main

import (
	"fmt"

	_ "github.com/pdrum/swagger-automation/docs" // This line is necessary for go-swagger to find your docs!
)

func main() {
	server := server{}
	if !server.serverStarted {
		fmt.Println("Welcome to github-stars server :)) .... ")
		server.startServer()
	}
}
