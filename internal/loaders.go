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

	return write(f, path)
}

func write(f *hclwrite.File, path string) error {
	tfFile, err := os.Create(path)
	if err != nil {
		return err
	}

	if _, err = tfFile.Write(f.Bytes()); err != nil {
		return err
	}

	return nil
}
