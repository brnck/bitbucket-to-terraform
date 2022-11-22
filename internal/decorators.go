package internal

import "github.com/ktrysmt/go-bitbucket"

type bitbucketRepositoryDecorator struct {
	branchRestriction []*bitbucket.BranchRestrictions
	pipelineConfig    *bitbucket.Pipeline
	bitbucket.Repository
}
