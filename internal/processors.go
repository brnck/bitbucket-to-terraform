package internal

import (
	"encoding/json"
	"fmt"
	"github.com/brnck/bitbucket-to-terraform/internal/config"
	"github.com/brnck/bitbucket-to-terraform/pkg/utils"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/ktrysmt/go-bitbucket"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
)

type Repository struct {
	config          *config.Config
	bitbucketClient *bitbucket.Client
}

func NewProcessors(c *config.Config, bbClient *bitbucket.Client) *Repository {
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

	log.Infof("Found %d projects", len(p.Items))
	log.Debugln("Transforming items to Terraform blocks")

	var projects []*hclwrite.Block
	var importStatements []string
	for _, v := range p.Items {
		block := transformToProjectBlock(&v)
		projects = append(projects, block)

		importStatement := transformToTerraformImportStatement(
			"bitbucket_project",
			utils.TransformStringToBeTFCompliant(v.Name),
			fmt.Sprintf("%s/%s", r.config.BitbucketWorkspace, v.Key),
		)
		importStatements = append(importStatements, importStatement)
	}

	if r.config.GenerateImportStatements {
		log.Infoln("Generate import statements flag is set to true. Will generate shell script")
		if err := writeTerraformImportStatementsToFile(
			importStatements,
			fmt.Sprintf("%s/%s.sh", r.config.ImportStatementsPath, "projects"),
		); err != nil {
			log.WithError(err)
			return err
		}
	}

	return writeTerraformBlocksToFile(projects, fmt.Sprintf("%s/%s.tf", r.config.Projects.Path, "projects"))
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

	log.Infof("Found %d repositories", len(res.Items))

	var wg sync.WaitGroup
	ch := make(chan bitbucketRepositoryDecorator, len(res.Items))
	wg.Add(len(res.Items))
	for _, v := range res.Items {
		go r.getInformationAboutRepository(v, ch, &wg)
	}
	wg.Wait()
	close(ch)

	log.Debugln("Transforming repositories to Terraform blocks")
	projects := make(map[string][]*hclwrite.Block)
	importStatements := make(map[string][]string)
	for brd := range ch {
		block := transformToRepositoryModuleBlock(&brd)
		name := utils.TransformStringToBeTFCompliant(brd.Project.Name)
		projects[name] = append(projects[name], block)

		importStatements[name] = append(
			importStatements[name],
			transformToTerraformModuleImportStatements(&brd, r.config.BitbucketWorkspace)...,
		)
	}

	log.Infoln("Writing repositories to files")
	if r.config.GenerateImportStatements {
		log.Infoln("Generate import statements flag is set to true. Will generate shell script")
		for index, statements := range importStatements {
			if err := writeTerraformImportStatementsToFile(
				statements,
				fmt.Sprintf("%s/project-%s-repositories.sh", r.config.ImportStatementsPath, index),
			); err != nil {
				log.WithError(err)
				return err
			}
		}
	}

	for name, blocks := range projects {
		if err := writeTerraformBlocksToFile(blocks, fmt.Sprintf(
			"%s/project-%s-repositories.tf",
			r.config.Repositories.Path,
			name,
		)); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) getInformationAboutRepository(
	v bitbucket.Repository,
	out chan<- bitbucketRepositoryDecorator,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	brd := bitbucketRepositoryDecorator{Repository: v}
	log.Infof("Fetching branch restrictions for repository %s", v.Name)
	br, err := r.getBranchRestrictions(v.Slug)
	if err != nil {
		if err.Error() == "404 Not Found" {
			log.Warnf("Branch restrictions for repo \"%s\" not found", v.Name)
		} else {
			log.Errorf("Error while fetching branch restriction for repository %s", v.Name)
		}
	} else {
		log.Infof("Branch restrictions for repository %s are fetched", v.Name)
	}

	brd.branchRestriction = br

	log.Infof("Fetching pipeline configuration for repository %s", v.Name)
	pp, err := r.getPipelineConfiguration(v.Slug)
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			log.Debugf("Branch restrictions for repo \"%s\" not found", v.Name)
		} else {
			log.Errorf("Error while fetching pipeline for repository %s", v.Name)
		}
	} else {
		log.Infof("Pipeline configuration for repository %s is fetched", v.Name)
	}
	brd.pipelineConfig = pp

	out <- brd
}

func (r *Repository) getBranchRestrictions(slug string) ([]*bitbucket.BranchRestrictions, error) {
	var branchRestrictions []*bitbucket.BranchRestrictions

	br, err := r.bitbucketClient.Repositories.BranchRestrictions.Gets(&bitbucket.BranchRestrictionsOptions{
		Owner:    r.config.BitbucketWorkspace,
		RepoSlug: slug,
	})

	if err != nil {
		return nil, err
	}

	brj, err := json.Marshal(br.(map[string]interface{})["values"])
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(brj, &branchRestrictions); err != nil {
		return nil, err
	}

	return branchRestrictions, nil
}

func (r *Repository) getPipelineConfiguration(slug string) (*bitbucket.Pipeline, error) {
	pc, err := r.bitbucketClient.Repositories.Repository.GetPipelineConfig(&bitbucket.RepositoryPipelineOptions{
		Owner:    r.config.BitbucketWorkspace,
		RepoSlug: slug,
	})
	if err != nil {
		return nil, err
	}

	return pc, nil
}
