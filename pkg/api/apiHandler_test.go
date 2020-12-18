package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/90lantran/github-star/internal/constants"
	"github.com/90lantran/github-star/internal/model"
)

const (
	healthExpectedResponse = "{\"message\":\"the server is up!\"}"
	invalidRequestResponse = "{\"error\":\"invalid request\"}"
	getStarsResponse       = "{\"totalStars\":19,\"invalidRepos\":[\"tingo-org/homebrew-tools\",\"tiygo-org/tinyfont\",\"tinygo-org/tinyfnt\"],\"validRepos\":{\"tinygo-org/tinyfont\":19}}"
	internalServerResponse = "{\"error\":\"cannot connect to github\"}"
)

type Route struct {
	http.Handler
	Method string
	Path   string
}

func (r *Route) Test(w http.ResponseWriter, req *http.Request) {
	m := mux.NewRouter()
	m.Handle(r.Path, r).Methods(r.Method)
	m.ServeHTTP(w, req)
}

func getStarsTestRoute() *Route {
	return &Route{
		Method:  "POST",
		Path:    constants.APIGetStarsEndpoint,
		Handler: http.HandlerFunc(GetStars),
	}
}
func cleanResponse(responseBody string) string {
	return strings.Replace(responseBody, "\n", "", 1)
}

func TestHealth(t *testing.T) {
	Convey("Given a health check request send to "+constants.APIHeathEndpoint, t, func() {
		request := httptest.NewRequest("GET", constants.APIHeathEndpoint, nil)
		response := httptest.NewRecorder()
		route := &Route{
			Method:  "GET",
			Path:    constants.APIHeathEndpoint,
			Handler: http.HandlerFunc(Health),
		}
		Convey("When the request is handled by the router", func() {
			route.Test(response, request)
			Convey("Then we should get expected reponse and success http response code", func() {
				So(response.Code, ShouldEqual, 200)
				So(string(cleanResponse(response.Body.String())), ShouldEqual, healthExpectedResponse)
			})
		})

	})
}

func TestGetStars(t *testing.T) {
	Convey("Given an invalid request send to "+constants.APIGetStarsEndpoint, t, func() {
		body := []byte(`
				{}
		`)
		request := httptest.NewRequest("POST", constants.APIGetStarsEndpoint, bytes.NewReader(body))
		response := httptest.NewRecorder()
		route := getStarsTestRoute()
		Convey("When the request is handled by the router", func() {
			route.Test(response, request)
			Convey("Then we should get response with bad http response code", func() {
				So(response.Code, ShouldEqual, 400)
				So(string(cleanResponse(response.Body.String())), ShouldEqual, invalidRequestResponse)
			})
		})
	})

	Convey("Given a valid request send to "+constants.APIGetStarsEndpoint, t, func() {
		body := []byte(`
				{"input":["tingo-org/homebrew-tools","tinygo-org/tinyfont","tiygo-org/tinyfont","tinygo-org/tinyfnt"]}
		`)
		request := httptest.NewRequest("POST", constants.APIGetStarsEndpoint, bytes.NewReader(body))
		response := httptest.NewRecorder()
		route := getStarsTestRoute()
		Convey("When the request is handled by the router", func() {
			route.Test(response, request)
			Convey("Then we should get response with success http response code", func() {
				So(response.Code, ShouldEqual, 200)
				So(string(cleanResponse(response.Body.String())), ShouldEqual, getStarsResponse)
			})
		})
	})

	Convey("Given a request send to "+constants.APIGetStarsEndpoint+"and cannot connet to github", t, func() {
		body := []byte(`
				{"input":["tingo-org/homebrew-tools","tinygo-org/tinyfont","tiygo-org/tinyfont","tinygo-org/tinyfnt"]}
		`)
		setGitService(model.GithubService{})
		request := httptest.NewRequest("POST", constants.APIGetStarsEndpoint, bytes.NewReader(body))
		response := httptest.NewRecorder()
		route := getStarsTestRoute()
		Convey("When the request is handled by the router", func() {
			route.Test(response, request)
			Convey("Then we should get response with internal http response code", func() {
				So(response.Code, ShouldEqual, 500)
				So(string(cleanResponse(response.Body.String())), ShouldEqual, internalServerResponse)
			})
		})
	})
}
