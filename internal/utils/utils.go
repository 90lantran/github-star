package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/90lantran/github-star/internal/model"
)

func ValidateInput(input string) error {
	var validInput = regexp.MustCompile(`^[a-zA-Z0-9\_\-\.]+\/[a-zA-Z0-9\_\-\.]+$`)

	if !validInput.MatchString(strings.TrimSpace(input)) {
		return fmt.Errorf("input list %s is not valid. Valid format is list of organization/repository", input)
	}

	return nil
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func ListAllReposForAnOrg(gitService model.GithubService, orgName string) ([]*github.Repository, error) {
	var allRepos []*github.Repository

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "99a7fae4ff499fb007b35ea8b34261d6c2c7d7de"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, orgName, gitService.Opt)
		//repos, resp, err := gitService.Client.Repositories.ListByOrg(gitService.Ctx, orgName, gitService.Opt)
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
