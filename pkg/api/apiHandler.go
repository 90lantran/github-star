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

const (
	success              = "success"
	failure              = "failure"
	inputNotValidMessage = "At least one of the input is not valid"
)

func init() {
	gitService = model.GithubService{
		Ctx:    context.Background(),
		Client: github.NewClient(nil),
		Opt: &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{PerPage: 100}},
	}
}

func setGitService(mock model.GithubService) {
	gitService = mock
}

// Base is just for minikube deployment
func Base(w http.ResponseWriter, r *http.Request) {
	response := model.HealthResponse{
		Message: constants.MinikubeMessage,
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}

// Health is bussiness logic for /health endpoint
func Health(w http.ResponseWriter, r *http.Request) {
	response := model.HealthResponse{
		Message: constants.HealthCheckResponse,
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}

// GetStars is bussiness logic for /get-stars endpoint
func GetStars(w http.ResponseWriter, r *http.Request) {
	// parse payload
	var req model.Request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		log.Printf("... cannot decode request %v\n", err)
		utils.RespondWithJSON(w, http.StatusBadRequest, model.Response{Error: err.Error(), Status: failure})
		return
	}
	defer r.Body.Close()

	if req.Input == nil {
		log.Printf("... invalid request. Must contain \"input:\" %v", req)
		utils.RespondWithJSON(w, http.StatusBadRequest, model.Response{Error: "invalid request. Must contain 'input:'", Status: failure})
		return
	}

	// start couting stars
	var totalCount int64
	validRepos := make([]model.MapNameStar, 0)
	invalidRepos := make([]string, 0)
	// caching, save number of call to github API
	seenOrgs := make(map[string][]*github.Repository)

	if gitService.Client == nil || gitService.Ctx == nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, model.Response{Error: "cannot connect to github", Status: failure})
		return
	}

	for _, input := range *req.Input {
		log.Printf("... check if each element in the format organization/repository")
		if err := utils.ValidateInput(input); err != nil {
			log.Printf("%s is not in right format orginazation/repo\n", input)
			invalidRepos = append(invalidRepos, input)
			continue
		}
		log.Printf("...processing input %v\n", input)
		token := strings.Split(input, "/")

		allRepos, ok := seenOrgs[token[0]]
		if !ok {
			log.Printf("%s is not seen before\n", token[0])
			results, err := utils.ListAllReposForAnOrg(gitService, token[0])
			if err == nil {
				allRepos = results
				seenOrgs[token[0]] = results
				log.Printf("...added all repos of %s to cache\n", token[0])
			} else {
				log.Printf("%s is not a valid org\n", token[0])
				invalidRepos = append(invalidRepos, input)
				continue
			}
		}

		if count := utils.GetStarsForRepo(allRepos, token[1]); count != -1 {
			validRepos = append(validRepos, model.MapNameStar{input, count})
			totalCount += count
		} else {
			log.Printf("%s is not a valid repo in the organization %s \n", token[1], token[0])
			invalidRepos = append(invalidRepos, input)
		}

	}
	payload := model.Payload{
		TotalStars:   totalCount,
		ValidRepos:   validRepos,
		InvalidRepos: invalidRepos,
	}
	resp := model.Response{
		Pl:     &payload,
		Status: success,
	}

	if len(invalidRepos) != 0 {
		resp.Error = inputNotValidMessage
	}
	log.Println("finished request")
	log.Printf("Response: %+v\n", resp)
	utils.RespondWithJSON(w, http.StatusOK, resp)
}
