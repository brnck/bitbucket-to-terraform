package main

import (
	"flag"
	"fmt"
	"github.com/brnck/bitbucket-to-terraform/internal/config"
	"github.com/brnck/bitbucket-to-terraform/internal/handlers"
	"github.com/ktrysmt/go-bitbucket"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	var cfg config.Config

	initGlobalFlags(&cfg)
	initLogger(&cfg)
	initGenericFlags(&cfg)
	initResourceFlags(&cfg.Projects, "projects")
	initResourceFlags(&cfg.Repositories, "repositories")
	flag.Parse()

	log.Debugln("Initializing Bitbucket client")
	bbClient := bitbucket.NewBasicAuth(cfg.BitbucketUsername, cfg.BitbucketPassword)
	handler := handlers.New(&cfg, bbClient)

	log.Debugln("Finished initializing Bitbucket client")
	log.Infof("Selected bitbucket workspace: %s", cfg.BitbucketWorkspace)
	log.Infof("User that will be used for communication with bitbucket: %s", cfg.BitbucketUsername)
	if cfg.Projects.Fetch {
		log.Infoln("Projects fetch flag is set true. Processing projects")
		if err := handler.ProcessProjects(); err != nil {
			log.Fatal("error occurred while processing projects: ", err)
		}
	}

	if cfg.Repositories.Fetch {
		log.Infoln("Repositories fetch flag is set true. Processing repositories")
		if err := handler.ProcessRepositories(); err != nil {
			log.Fatal("error occurred while processing repositories: ", err)
		}
	}
}

func initGlobalFlags(c *config.Config) {
	flag.BoolVar(
		&c.DryRun,
		"dry-run",
		false,
		"Simulate data extraction, transformation and loading to specified paths",
	)

	flag.UintVar(&c.LogLevel, "verbose", 7, "Log severity level [1-7]")
}

func initLogger(c *config.Config) {
	if c.LogLevel > 7 {
		c.LogLevel = 7
	} else if c.LogLevel < 1 {
		c.LogLevel = 1
	}

	log.SetOutput(os.Stdout)
	level := log.AllLevels[c.LogLevel-1]
	log.SetLevel(level)
}

func initGenericFlags(c *config.Config) {
	flag.StringVar(&c.BitbucketUsername, "bitbucket-username", "username", "Bitbucket username")
	flag.StringVar(&c.BitbucketPassword, "bitbucket-password", "password", "Bitbucket password")
	flag.StringVar(&c.BitbucketWorkspace, "bitbucket-workspace", c.BitbucketUsername, "Which workspace to use")
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
		fmt.Sprintf("Where to extract %s", resourceName),
	)
	flag.BoolVar(
		&r.SplitToFiles,
		fmt.Sprintf("split-%s-to-files", resourceName),
		false,
		"Should each resource be separate TF file",
	)
}
