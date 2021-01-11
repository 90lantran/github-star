package model

import (
	"context"

	"github.com/google/go-github/github"
)

type Request struct {
	Input *[]string `json:"input" validate:"required"`
}

type Response struct {
	Pl     *Payload `json:"payload,omitempty"`
	Error  []string `json:"error,omitempty"`
	Status string   `json:"status" validate:"required"`
}
type Payload struct {
	TotalStars   int64         `json:"totalStars,omitempty"`
	InvalidRepos []string      `json:"invalidRepos,omitempty"`
	ValidRepos   []MapNameStar `json:"validRepos,omitempty"`
}

type MapNameStar struct {
	Name string `json:"name"`
	Star int64  `json:"star(s)"`
}

type GithubService struct {
	Ctx    context.Context
	Client *github.Client
	Opt    *github.RepositoryListByOrgOptions
}

type HealthResponse struct {
	Message string `json:"message"`
}
