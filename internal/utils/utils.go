package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/go-github/github"

	"github.com/90lantran/github-star/internal/model"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func ListAllReposForAnOrg(gitService model.GithubService, orgName string) ([]*github.Repository, error) {
	var allRepos []*github.Repository
	for {
		repos, resp, err := gitService.Client.Repositories.ListByOrg(gitService.Ctx, orgName, gitService.Opt)
		if err != nil {
			fmt.Printf("the organization %s does not exist %v\n", orgName, err)
			return nil, errors.New("the organization does not exist")
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		gitService.Opt.Page = resp.NextPage
	}
	return allRepos, nil

}

func GetStarsForRepo(allRepos []*github.Repository, repoName string) int {
	//fmt.Printf("repos %v", allRepos)

	for _, repo := range allRepos {
		if repo.GetName() == repoName {
			count := repo.GetStargazersCount()
			fmt.Printf("Found: %s\n", repoName)
			return count
		}
	}
	return -1
}
