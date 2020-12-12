package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

func getStars(w http.ResponseWriter, r *http.Request) {
	// b082d2cd0e5e4202f31a
	// 76b2b47503154546fc0393b70c3b488d3f6d66a1

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "08dcd33b87acba14d8630efbf2ae2736d885ad53"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 30},
	}
	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, "tinygo-org", opt)
		if err != nil {
			fmt.Printf("The organization does not exist %v", err)
			return
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	// client.Repositories
	// list all organizations for user "willnorris"
	// repos, _, err := client.Repositories.List(context.Background(), "", nil)

	// if err != nil {
	// 	fmt.Printf("Cannot call go-github %v", err)
	// }
	for _, repo := range allRepos {
		if repo.GetName() == "llvm-project" {
			fmt.Printf("Number of star(s): %d\n", repo.GetStargazersCount())
		}

	}

	fmt.Fprintf(w, "Done")
}

func main() {
	fmt.Println("... start server")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/getStars", getStars)
	log.Fatal(http.ListenAndServe(":8080", router))
}
