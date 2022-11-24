# bitbucket-to-terraform

Tool for importing resources from Bitbucket to Terraform using BitbucketAPI

## Context

This CLI app can fetch certain resources using Bitbucket API, iterate through the list of those resources and convert
them Terraform resources while, also, providing an ability to generate `terraform import ...` statements to ensure smooth
and easy terraforming of the Bitbucket environment.

By design this app prepares TF resources for [zahiar/bitbucket](https://registry.terraform.io/providers/zahiar/bitbucket/1.3.0)
terraform provider

## Usage

```shell
~$: bb2tf --help
Usage of bb2tf:
  -bitbucket-password string
    	Bitbucket password (default "password")
  -bitbucket-username string
    	Bitbucket username (default "username")
  -bitbucket-workspace string
    	Which workspace to use (default "username")
  -fetch-projects
    	Fetch projects from the Bitbucket (default true)
  -fetch-repositories
    	Fetch repositories from the Bitbucket (default true)
  -generate-import-statements terraform import <...>
    	Generates shell script terraform import <...> (default true)
  -import-statements-path string
    	If --generate-import-statements=true, it will be used as path for the file (default "./")
  -load-projects-path string
    	Where to extract projects (folder path, not file) (default "./")
  -load-repositories-path string
    	Where to extract repositories (folder path, not file) (default "./")
  -verbose int
    	Log severity level [1-7] (default 4)
```
