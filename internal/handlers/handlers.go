package handlers

import (
	"github.com/brnck/bitbucket-to-terraform/internal/config"
	"github.com/ktrysmt/go-bitbucket"
)

type Repository struct {
	config          *config.Config
	bitbucketClient *bitbucket.Client
}

func New(c *config.Config, bbClient *bitbucket.Client) *Repository {
	return &Repository{
		config:          c,
		bitbucketClient: bbClient,
	}
}

func (r *Repository) ProcessProjects() error {
	_, err := r.bitbucketClient.Workspaces.Projects(r.config.BitbucketWorkspace)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ProcessRepositories() error {
	return nil
}
