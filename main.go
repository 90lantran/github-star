package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

type request struct {
	Input *[]string `json:"input"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func getStars(w http.ResponseWriter, r *http.Request) {
	// b082d2cd0e5e4202f31a
	// 76b2b47503154546fc0393b70c3b488d3f6d66a1

	// should be in app.go
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "08dcd33b87acba14d8630efbf2ae2736d885ad53"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 30},
	}

	// parse payload
	var p request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	log.Printf("request: %v", *p.Input)
	count := 0
	for _, input := range *p.Input {
		fmt.Printf("Processing input %v\n", input)
		token := strings.Split(input, "/")
		count += getStarsForEachRepo(token[0], token[1], client, ctx, opt)
	}
	// get all pages of results
	// var count int
	// var allRepos []*github.Repository
	// for {
	// 	repos, resp, err := client.Repositories.ListByOrg(ctx, "tinygo-org", opt)
	// 	if err != nil {
	// 		fmt.Printf("The organization does not exist %v", err)
	// 		return
	// 	}
	// 	allRepos = append(allRepos, repos...)
	// 	if resp.NextPage == 0 {
	// 		break
	// 	}
	// 	opt.Page = resp.NextPage
	// }

	// for _, repo := range allRepos {
	// 	if repo.GetName() == "llvm-project" {
	// 		count += repo.GetStargazersCount()
	// 		//fmt.Printf("Number of star(s): %d\n", repo.GetStargazersCount())
	// 	}
	// }
	log.Println(count)
	respondWithJSON(w, http.StatusOK, count)
}

func getStarsForEachRepo(orgName string, repoName string, client *github.Client, ctx context.Context, opt *github.RepositoryListByOrgOptions) int {
	var count int
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, orgName, opt)
		if err != nil {
			fmt.Printf("The organization does not exist %v", err)
			continue
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, repo := range allRepos {
		if repo.GetName() == repoName {
			count += repo.GetStargazersCount()
			fmt.Printf("Found: %s\n", repoName)
		}
	}
	return count
}

func main() {
	fmt.Println("... start server")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/get-stars", getStars).Methods("POST")
	//router.HandleFunc("/getStars", getStars)
	log.Fatal(http.ListenAndServe(":8080", router))
}
