package utils

import (
	"encoding/json"
	"net/http"

	"github.com/google/go-github/github"

	"github.com/90lantran/github-star/internal/model"
)

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
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		gitService.Opt.Page = resp.NextPage
	}
	return allRepos, nil

}

func GetStarsForRepo(allRepos []*github.Repository, repoName string) int64 {
	for _, repo := range allRepos {
		if repo.GetName() == repoName {
			count := repo.GetStargazersCount()
			return int64(count)
		}
	}
	return -1
}
