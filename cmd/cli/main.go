package main

import (
	"flag"
	"fmt"
	"github.com/brnck/bitbucket-to-terraform/internal"
	"github.com/brnck/bitbucket-to-terraform/internal/config"
	"github.com/ktrysmt/go-bitbucket"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	var cfg config.Config

	flag.IntVar(&cfg.LogLevel, "verbose", 4, "Log severity level [1-7]")
	flag.StringVar(&cfg.BitbucketUsername, "bitbucket-username", "username", "Bitbucket username")
	flag.StringVar(&cfg.BitbucketPassword, "bitbucket-password", "password", "Bitbucket password")
	flag.StringVar(&cfg.BitbucketWorkspace, "bitbucket-workspace", cfg.BitbucketUsername, "Which workspace to use")

	flag.BoolVar(
		&cfg.GenerateImportStatements,
		"generate-import-statements",
		true,
		"Generates shell script `terraform import <...>`",
	)
	flag.StringVar(
		&cfg.ImportStatementsPath,
		"import-statements-path",
		"./",
		"If --generate-import-statements=true, it will be used as path for the file",
	)

	initResourceFlags(&cfg.Projects, "projects")
	initResourceFlags(&cfg.Repositories, "repositories")
	flag.Parse()

	initLogger(&cfg)
	log.Debugln("Initializing Bitbucket client")

	bbClient := bitbucket.NewBasicAuth(cfg.BitbucketUsername, cfg.BitbucketPassword)
	bbClient.Pagelen = 100
	bbClient.DisableAutoPaging = false
	processor := internal.NewProcessors(&cfg, bbClient)

	log.Debugln("Finished initializing Bitbucket client")
	log.Infof("Selected bitbucket workspace: %s", cfg.BitbucketWorkspace)
	log.Infof("User that will be used for communication with bitbucket: %s", cfg.BitbucketUsername)

	if cfg.Projects.Fetch {
		log.Infoln("Projects fetch flag is set true. Processing projects")
		if err := processor.ProcessProjects(); err != nil {
			log.Fatal("error occurred while processing projects: ", err)
		}
	}

	if cfg.Repositories.Fetch {
		log.Infoln("Repositories fetch flag is set true. Processing repositories")
		if err := processor.ProcessRepositories(); err != nil {
			log.Fatal("error occurred while processing repositories: ", err)
		}
	}
}

func initLogger(c *config.Config) {
	if c.LogLevel > 7 {
		c.LogLevel = 7
	} else if c.LogLevel < 1 {
		c.LogLevel = 1
	}

	log.SetOutput(os.Stdout)
	log.SetLevel(log.AllLevels[c.LogLevel-1])
}

func initResourceFlags(r *config.ResourceFetchConfig, resourceName string) {
	flag.BoolVar(
		&r.Fetch,
		fmt.Sprintf("fetch-%s", resourceName),
		true,
		fmt.Sprintf("Fetch %s from the Bitbucket", resourceName),
	)
	flag.StringVar(
		&r.Path,
		fmt.Sprintf("load-%s-path", resourceName),
		"./",
		fmt.Sprintf("Where to extract %s (folder path, not file)", resourceName),
	)
}
