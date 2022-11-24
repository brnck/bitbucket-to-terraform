package internal

import (
	"fmt"
	"github.com/brnck/bitbucket-to-terraform/pkg/utils"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/ktrysmt/go-bitbucket"
	"github.com/zclconf/go-cty/cty"
)

var restrictionsToSkip = []string{
	"push",
	"force",
	"delete",
	"restrict_merges",
}

var restrictionsWithFlowControl = []string{
	"require_no_changes_requested",
	"require_tasks_to_be_completed",
	"require_passing_builds_to_merge",
	"enforce_merge_checks",
	"allow_auto_merge_when_builds_pass",
	"reset_pullrequest_approvals_on_change",
	"reset_pullrequest_changes_requested_on_change",
	"smart_reset_pullrequest_approvals",
}

// transformToProjectBlock creates terraform resource from Bitbucket's project struct
func transformToProjectBlock(p *bitbucket.Project) *hclwrite.Block {
	block := hclwrite.NewBlock("resource", []string{
		"bitbucket_project",
		utils.TransformStringToBeTFCompliant(p.Name),
	})

	body := block.Body()

	workspaces := hclwrite.Tokens{
		{Type: hclsyntax.TokenIdent, Bytes: []byte(`var.workspace`)},
	}
	body.SetAttributeRaw("workspace", workspaces)

	body.SetAttributeValue("name", cty.StringVal(p.Name))
	body.SetAttributeValue("key", cty.StringVal(p.Key))
	body.SetAttributeValue("description", cty.StringVal(p.Description))
	body.SetAttributeValue("is_private", cty.BoolVal(p.Is_private))

	return block
}

// transformToRepositoryModuleBlock creates terraform resource from Bitbucket's project struct
func transformToRepositoryModuleBlock(r *bitbucketRepositoryDecorator) *hclwrite.Block {
	block := hclwrite.NewBlock(
		"module",
		[]string{utils.TransformStringToBeTFCompliant(r.Slug)},
	)

	body := block.Body()

	body.SetAttributeValue("source", cty.StringVal("../modules/repository"))
	body.SetAttributeValue("name", cty.StringVal(r.Slug))
	body.SetAttributeValue("description", cty.StringVal(r.Description))
	body.SetAttributeValue("fork_policy", cty.StringVal(r.Fork_policy))
	body.SetAttributeValue("default_branch", cty.StringVal(r.Mainbranch.Name))

	workspaces := hclwrite.Tokens{
		{Type: hclsyntax.TokenIdent, Bytes: []byte(`var.workspace`)},
	}
	body.SetAttributeRaw("workspace_name", workspaces)

	projectKey := hclwrite.Tokens{
		{
			Type: hclsyntax.TokenIdent,
			Bytes: []byte(fmt.Sprintf(
				"bitbucket_project.%s.key",
				utils.TransformStringToBeTFCompliant(r.Project.Name),
			)),
		},
	}
	body.SetAttributeRaw("project_key", projectKey)

	if r.pipelineConfig == nil {
		body.SetAttributeValue("pipelines_enabled", cty.BoolVal(false))
	} else {
		body.SetAttributeValue("pipelines_enabled", cty.BoolVal(r.pipelineConfig.Enabled))
	}

	for _, restriction := range r.branchRestriction {
		if utils.StringExistsInList(restrictionsToSkip, restriction.Kind) ||
			r.Mainbranch.Name != restriction.Pattern {
			continue
		}

		if utils.StringExistsInList(restrictionsWithFlowControl, restriction.Kind) {
			body.SetAttributeValue(restriction.Kind, cty.BoolVal(true))
		}

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

func transformToTerraformModuleImportStatements(r *bitbucketRepositoryDecorator, workspace string) []string {
	var statements []string

	repoImport := transformToTerraformImportStatement(
		"module",
		fmt.Sprintf("%s.bitbucket_repository.this", utils.TransformStringToBeTFCompliant(r.Name)),
		fmt.Sprintf("%s/%s", workspace, r.Slug),
	)
	statements = append(statements, repoImport)

	branchRestriction := fmt.Sprintf("%s.bitbucket_branch_restriction.", utils.TransformStringToBeTFCompliant(r.Name))

	for _, restriction := range r.branchRestriction {
		if utils.StringExistsInList(restrictionsToSkip, restriction.Kind) ||
			r.Mainbranch.Name != restriction.Pattern {
			continue
		}

		br := branchRestriction + restriction.Kind + "\\[0\\]"

		brts := transformToTerraformImportStatement("module", br,
			fmt.Sprintf("%s/%s/%d", workspace, r.Slug, restriction.ID),
		)

		statements = append(statements, brts)
	}

	return statements
}

func transformToTerraformImportStatement(tfResource, tfResourceName, bitbucketResource string) string {
	return fmt.Sprintf("terraform import %s.%s %s", tfResource, tfResourceName, bitbucketResource)
}
