package main

import (
	"log"
	"net/http"

	"github.com/90lantran/github-star/pkg/route"
)

// Server represents the server side component of the API
type Server struct {
	serverStarted bool
	// DB connection  if we want to extend
}

// StartServer starts the server
func (s *Server) StartServer() {
	s.serverStarted = true
	router := route.GetRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
