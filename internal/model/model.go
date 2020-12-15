package model

import (
	"context"

	"github.com/google/go-github/github"
)

type Request struct {
	Input *[]string `json:"input"`
}

type Response struct {
	TotalStars   int            `json:"totalStars,omitempty"`
	InvalidRepos []string       `json:"invalidRepos,omitempty"`
	ValidRepos   map[string]int `json:"validRepos,omitempty"`
}

type GithubService struct {
	Ctx    context.Context
	Client *github.Client
	Opt    *github.RepositoryListByOrgOptions
}
