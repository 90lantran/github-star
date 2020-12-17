package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-github/github"

	"github.com/90lantran/github-star/internal/constants"
	"github.com/90lantran/github-star/internal/model"
	"github.com/90lantran/github-star/internal/utils"
)

var gitService model.GithubService

func init() {
	gitService = model.GithubService{
		Ctx:    context.Background(),
		Client: github.NewClient(nil),
		Opt: &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{PerPage: 30}},
	}
}

// Health is bussiness logic for /health endpoint
func Health(w http.ResponseWriter, r *http.Request) {
	resp := constants.HealthCheckResponse
	utils.RespondWithJSON(w, http.StatusOK, resp)
}

// GetStars is bussiness logic for /get-stars endpoint
func GetStars(w http.ResponseWriter, r *http.Request) {
	// gitService := model.GithubService{
	// 	Ctx:    context.Background(),
	// 	Client: github.NewClient(nil),
	// 	Opt: &github.RepositoryListByOrgOptions{
	// 		ListOptions: github.ListOptions{PerPage: 30}},
	// }
	// b082d2cd0e5e4202f31a
	// 76b2b47503154546fc0393b70c3b488d3f6d66a1

	// gitService := model.GithubService{
	// 	Ctx:    context.Background(),
	// 	Client: github.NewClient(nil),
	// 	Opt:    &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{PerPage: 30}},
	// }
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
	var req model.Request
	decoder := json.NewDecoder(r.Body)
	var err error
	if err = decoder.Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot decode the payload")
		log.Fatalf("internal server error %v\n", err)
	}
	defer r.Body.Close()

	if req.Input == nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// start couting
	var totalCount int64
	var resp model.Response
	validRepos := make(map[string]int64)
	invalidRepos := make([]string, 0)
	// caching, save number of to github API
	seenOrgs := make(map[string][]*github.Repository)

	for _, input := range *req.Input {
		log.Printf("...processing input %v\n", input)
		token := strings.Split(input, "/")
		allRepos, ok := seenOrgs[token[0]]
		if !ok {
			log.Printf("%s is not seen before\n", token[0])
			results, err := utils.ListAllReposForAnOrg(gitService, token[0])
			if err == nil {
				allRepos = results
				seenOrgs[token[0]] = results
				log.Printf("...added all repos of %s to map\n", token[0])
			} else {
				log.Printf("%s is not a valid org\n", token[0])
				invalidRepos = append(invalidRepos, input)
				continue
			}
		}

		if count := utils.GetStarsForRepo(allRepos, token[1]); count != -1 {
			validRepos[input] = count
			totalCount += count
		} else {
			log.Printf("%s is not a valid repo in the organization %s \n", token[1], token[0])
			invalidRepos = append(invalidRepos, input)
		}

	}
	resp.TotalStars = totalCount
	resp.ValidRepos = validRepos
	resp.InvalidRepos = invalidRepos
	log.Println("finished request")
	utils.RespondWithJSON(w, http.StatusOK, resp)
}
