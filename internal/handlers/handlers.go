package handlers

import (
	"fmt"
	"github.com/brnck/bitbucket-to-terraform/internal"
	"github.com/brnck/bitbucket-to-terraform/internal/config"
	"github.com/brnck/bitbucket-to-terraform/pkg/utils"
	"github.com/hashicorp/hcl/v2/hclwrite"
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

// ProcessProjects Get all Projects from the workspace, transform them to terraform blocks, write to file or files
func (r *Repository) ProcessProjects() error {
	p, err := r.bitbucketClient.Workspaces.Projects(r.config.BitbucketWorkspace)
	if err != nil {
		return err
	}

	projects := make(map[string]*hclwrite.Block)
	for _, v := range p.Items {
		block := internal.TransformToProjectBlock(&v, r.config.BitbucketWorkspace)
		name := utils.TransformStringToBeTFCompliant(v.Name)
		projects[name] = block
	}

	return writeToFile(&r.config.Projects, projects)
}

// ProcessRepositories Get all repositories from the workspace, transform them to terraform blocks, write to file or files
func (r *Repository) ProcessRepositories() error {
	opt := bitbucket.RepositoriesOptions{
		Owner: r.config.BitbucketWorkspace,
	}

	res, err := r.bitbucketClient.Repositories.ListForAccount(&opt)
	if err != nil {
		return err
	}

	repositories := make(map[string]*hclwrite.Block)
	for _, v := range res.Items {
		block := internal.TransformToRepositoryModuleBlock(&v, r.config.BitbucketWorkspace)
		name := utils.TransformStringToBeTFCompliant(v.Name)
		repositories[name] = block
	}

	return writeToFile(&r.config.Repositories, repositories)
}

func writeToFile(c *config.ResourceFetchConfig, resources map[string]*hclwrite.Block) error {
	if c.SplitToFiles {
		for name, block := range resources {
			if err := internal.WriteTerraformBlockToFile(block, fmt.Sprintf("%s/%s.tf", c.Path, name)); err != nil {
				return err
			}
		}

		return nil
	}

	return internal.WriteTerraformBlocksToFile(resources, fmt.Sprintf("%s/%s.tf", c.Path, "repositories"))
}
