# bitbucket-to-terraform

Tool for importing resources from Bitbucket to Terraform using BitbucketAPI

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
  -dry-run
    	Simulate data extraction, transformation and loading to specified paths
  -fetch-projects
    	Fetch projects from the Bitbucket (default true)
  -fetch-repositories
    	Fetch repositories from the Bitbucket (default true)
  -load-projects-path string
    	Where to extract projects (default "./")
  -load-repositories-path string
    	Where to extract repositories (default "./")
  -split-projects-to-files
    	Should each resource be separate TF file
  -split-repositories-to-files
    	Should each resource be separate TF file
  -verbose uint
    	Log severity level [1-7] (default 4)
```
