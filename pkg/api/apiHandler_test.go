package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/90lantran/github-star/internal/constants"
	"github.com/90lantran/github-star/internal/model"
)

const (
	healthExpectedResponse = "{\"message\":\"the server is up!\"}"
	emptyRequestResponse   = "{\"error\":[\"invalid request. Must contain 'input:'\"],\"status\":\"failure\"}"
	getStarsResponse       = "{\"payload\":{\"totalStars\":19,\"invalidRepos\":[\"tingo-org/homebrew-tools\",\"tiygo-org\",\"tinygo-org/tinyfnt\"],\"validRepos\":[{\"name\":\"tinygo-org/tinyfont\",\"star(s)\":19}]},\"error\":[\"GET https://api.github.com/orgs/tingo-org/repos?per_page=100: 404 Not Found []\",\"input list tiygo-org is not valid. Valid format is list of organization/repository\",\"tinyfnt is not a valid repo in the organization tinygo-org\"],\"status\":\"success\"}"
	internalServerResponse = "{\"error\":[\"cannot connect to github\"],\"status\":\"failure\"}"
	invalidRequestResponse = "{\"error\":[\"json: cannot unmarshal string into Go struct field Request.input of type []string\"],\"status\":\"failure\"}"
	ratelimitRespons       = "403 API rate limit"
	notRefeshPagination    = "{\"payload\":{\"totalStars\":874,\"invalidRepos\":[\"twilio/twilio-python\"],\"validRepos\":[{\"name\":\"twitter/rezolus\",\"star(s)\":874}]},\"error\":[\"twilio-python is not a valid repo in the organization twilio\"],\"status\":\"success\"}"
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
	Convey("Given an empty request send to "+constants.APIGetStarsEndpoint, t, func() {
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
				So(string(cleanResponse(response.Body.String())), ShouldEqual, emptyRequestResponse)
			})
		})
	})

	Convey("Given an invalid request send to "+constants.APIGetStarsEndpoint, t, func() {
		body := []byte(`
				{"input":"tinygo-org/tinygo-site"}
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
				{"input":["tingo-org/homebrew-tools","tinygo-org/tinyfont","tiygo-org","tinygo-org/tinyfnt"]}
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

	Convey("Given a request send to "+constants.APIGetStarsEndpoint+"and did not refresh pagination", t, func() {
		body := []byte(`
				{"input":["twitter/rezolus", "twilio/twilio-python"]}
		`)
		setFlag(false)
		request := httptest.NewRequest("POST", constants.APIGetStarsEndpoint, bytes.NewReader(body))
		response := httptest.NewRecorder()
		route := getStarsTestRoute()
		Convey("When the request is handled by the router", func() {

			route.Test(response, request)
			Convey("Then we should get response with 200 status and a valid input will be listed as invalid", func() {
				So(response.Code, ShouldEqual, 200)
				fmt.Println(string(cleanResponse(response.Body.String())))
				So(string(cleanResponse(response.Body.String())), ShouldContainSubstring, notRefeshPagination)
			})
		})
	})

	Convey("Given a request send to "+constants.APIGetStarsEndpoint+"and cannot connet to github", t, func() {
		body := []byte(`
				{"input":["tingo-org/homebrew-tools","tinygo-org/tinyfont","tiygo-org","tinygo-org/tinyfnt"]}
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

	Convey("Given a request send to "+constants.APIGetStarsEndpoint+"and rate limit is hit", t, func() {
		body := []byte(`
				{"input":["google/trax", "microsoft/TypeScript"]}
		`)
		setGitService(
			model.GithubService{
				Opt:    &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{PerPage: 1}},
				Ctx:    context.Background(),
				Client: github.NewClient(nil),
			},
		)
		request := httptest.NewRequest("POST", constants.APIGetStarsEndpoint, bytes.NewReader(body))
		response := httptest.NewRecorder()
		route := getStarsTestRoute()
		Convey("When the request is handled by the router", func() {

			route.Test(response, request)
			Convey("Then we should get response with 200 and error 403 rate limit hit", func() {
				So(response.Code, ShouldEqual, 200)
				So(string(cleanResponse(response.Body.String())), ShouldContainSubstring, ratelimitRespons)
			})
		})
	})
}
