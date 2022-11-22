package internal

import (
	"github.com/brnck/bitbucket-to-terraform/pkg/utils"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/ktrysmt/go-bitbucket"
	"github.com/zclconf/go-cty/cty"
)

// transformToProjectBlock creates terraform resource from Bitbucket's project struct
func transformToProjectBlock(p *bitbucket.Project, workspace string) *hclwrite.Block {
	block := hclwrite.NewBlock("resource", []string{
		"bitbucket_project",
		utils.TransformStringToBeTFCompliant(p.Name),
	})

	body := block.Body()

	body.SetAttributeValue("workspace", cty.StringVal(workspace))
	body.SetAttributeValue("name", cty.StringVal(p.Name))
	body.SetAttributeValue("key", cty.StringVal(p.Key))
	body.SetAttributeValue("description", cty.StringVal(p.Description))
	body.SetAttributeValue("is_private", cty.BoolVal(p.Is_private))

	return block
}

// transformToRepositoryModuleBlock creates terraform resource from Bitbucket's project struct
func transformToRepositoryModuleBlock(r *bitbucketRepositoryDecorator, workspace string) *hclwrite.Block {
	block := hclwrite.NewBlock(
		"module",
		[]string{utils.TransformStringToBeTFCompliant(r.Slug)},
	)

	body := block.Body()

	body.SetAttributeValue("source", cty.StringVal("../modules/repository"))
	body.SetAttributeValue("name", cty.StringVal(r.Name))
	body.SetAttributeValue("description", cty.StringVal(r.Description))
	body.SetAttributeValue("workspace_name", cty.StringVal(workspace))
	body.SetAttributeValue("project_key", cty.StringVal(r.Project.Key))

	if r.pipelineConfig == nil {
		body.SetAttributeValue("pipelines_enabled", cty.BoolVal(false))
	} else {
		body.SetAttributeValue("pipelines_enabled", cty.BoolVal(r.pipelineConfig.Enabled))
	}

	for _, restriction := range r.branchRestriction {
		if restriction.Kind == "require_approvals_to_merge" {
			body.SetAttributeValue("required_approvals_to_merge", cty.NumberIntVal(int64(*restriction.Value)))
		}

		if restriction.Kind == "require_default_reviewer_approvals_to_merge" {
			body.SetAttributeValue(
				"require_default_reviewer_approvals_to_merge",
				cty.NumberIntVal(int64(*restriction.Value)),
			)
		}
	}

	return block
}
