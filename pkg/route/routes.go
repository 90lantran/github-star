package route

import (
	"github.com/gorilla/mux"

	"github.com/90lantran/github-star/internal/constants"
	"github.com/90lantran/github-star/pkg/api"
)

var router *mux.Router

func init() {
	router = mux.NewRouter().StrictSlash(false)
}

// GetRouter configures router to handle requests coming to server
func GetRouter() *mux.Router {
	router.HandleFunc(constants.APIBaseEndpoint, api.Base).Methods("GET").Name("base")
	router.HandleFunc(constants.APIHeathEndpoint, api.Health).Methods("GET").Name("health")
	router.HandleFunc(constants.APIGetStarsEndpoint, api.GetStars).Methods("POST").Name("get-stars")
	return router
}
