package internal

import (
	"github.com/brnck/bitbucket-to-terraform/pkg/utils"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/ktrysmt/go-bitbucket"
	"github.com/zclconf/go-cty/cty"
)

// TransformToProjectBlock creates terraform resource from Bitbucket's project struct
func TransformToProjectBlock(p *bitbucket.Project, workspace string) *hclwrite.Block {
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

// TransformToRepositoryModuleBlock creates terraform resource from Bitbucket's project struct
func TransformToRepositoryModuleBlock(r *bitbucket.Repository, workspace string) *hclwrite.Block {
	block := hclwrite.NewBlock(
		"module",
		[]string{utils.TransformStringToBeTFCompliant(r.Name)},
	)

	body := block.Body()

	body.SetAttributeValue("source", cty.StringVal("../modules/repository"))
	body.SetAttributeValue("name", cty.StringVal(r.Name))
	body.SetAttributeValue("description", cty.StringVal(r.Description))
	body.SetAttributeValue("workspace_name", cty.StringVal(workspace))
	body.SetAttributeValue("project_key", cty.StringVal(r.Project.Key))

	//body.SetAttributeValue("pipelines_enabled", cty.BoolVal(true))
	//body.SetAttributeValue("required_approvals_to_merge", cty.NumberUIntVal(1))
	//body.SetAttributeValue("require_default_reviewer_approvals_to_merge", cty.NumberUIntVal(1))

	return block
}
