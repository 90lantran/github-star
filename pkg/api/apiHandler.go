package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-github/github"

	"github.com/90lantran/github-star/internal/model"
	"github.com/90lantran/github-star/internal/utils"
)

// Health is bussiness logic for /health endpoint
func Health(w http.ResponseWriter, r *http.Request) {

}

// GetStars is bussiness logic for /get-stars endpoint
func GetStars(w http.ResponseWriter, r *http.Request) {
	// b082d2cd0e5e4202f31a
	// 76b2b47503154546fc0393b70c3b488d3f6d66a1

	gitService := model.GithubService{
		Ctx:    context.Background(),
		Client: github.NewClient(nil),
		Opt:    &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{PerPage: 30}},
	}
	// should be in app.go
	//ctx := context.Background()
	// ts := oauth2.StaticTokenSource(
	// 	&oauth2.Token{AccessToken: "08dcd33b87acba14d8630efbf2ae2736d885ad53"},
	// )
	//tc := oauth2.NewClient(ctx, ts)

	//client := github.NewClient(tc)
	// client := github.NewClient(nil)

	// opt := &github.RepositoryListByOrgOptions{
	// 	ListOptions: github.ListOptions{PerPage: 30},
	// }

	// parse payload
	var p model.Request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	log.Printf("request: %v", *p.Input)

	totalCount := 0
	validRepos := make(map[string]int)
	invalidRepos := make([]string, 0)
	var resp model.Response
	seenOrgs := make(map[string][]*github.Repository)
	for _, input := range *p.Input {
		fmt.Printf("Processing input %v\n", input)
		token := strings.Split(input, "/")
		allRepos, ok := seenOrgs[token[0]]
		if !ok {
			fmt.Printf("%s is not in the map\n", token[0])
			results, err := utils.ListAllReposForAnOrg(gitService, token[0])
			if err == nil {
				fmt.Printf("Find all repos for %s\n", token[0])
				allRepos = results
				seenOrgs[token[0]] = results

			} else {
				// invalid org
				fmt.Printf("%s is not a valid org\n", token[0])
				invalidRepos = append(invalidRepos, input)
				continue
			}

		}

		if count := utils.GetStarsForRepo(allRepos, token[1]); count != -1 {
			validRepos[input] = count
			totalCount += count
		} else {
			// invalid repo
			fmt.Printf("%s is not a valid repo\n", token[1])
			invalidRepos = append(invalidRepos, input)
		}

	}

	resp.TotalStars = totalCount
	resp.ValidRepos = validRepos
	resp.InvalidRepos = invalidRepos
	utils.RespondWithJSON(w, http.StatusOK, resp)
}
