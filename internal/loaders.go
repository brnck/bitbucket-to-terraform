package internal

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"os"
)

// writeTerraformBlocksToFile writes multiple terraform resources to file
func writeTerraformBlocksToFile(b []*hclwrite.Block, path string) error {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	for _, block := range b {
		rootBody.AppendBlock(block)
		rootBody.AppendNewline()
	}

	return write(path, f.Bytes())
}

func writeTerraformImportStatementsToFile(s []string, path string) error {
	body := `#!/bin/bash

set -euxo pipefail
`

	for _, statement := range s {
		body = body + statement + "\n"
	}

	return write(path, []byte(body))
}

func write(path string, bytes []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if _, err = f.Write(bytes); err != nil {
		return err
	}

	return nil
}
