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

	flag.UintVar(&cfg.LogLevel, "verbose", 7, "Log severity level [1-7]")
	flag.StringVar(&cfg.BitbucketUsername, "bitbucket-username", "username", "Bitbucket username")
	flag.StringVar(&cfg.BitbucketPassword, "bitbucket-password", "password", "Bitbucket password")
	flag.StringVar(&cfg.BitbucketWorkspace, "bitbucket-workspace", cfg.BitbucketUsername, "Which workspace to use")

	initLogger(&cfg)
	initResourceFlags(&cfg.Projects, "projects")
	initResourceFlags(&cfg.Repositories, "repositories")
	flag.Parse()

	log.Debugln("Initializing Bitbucket client")
	bbClient := bitbucket.NewBasicAuth(cfg.BitbucketUsername, cfg.BitbucketPassword)
	bbClient.Pagelen = 50
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
	flag.BoolVar(
		&r.SplitToFiles,
		fmt.Sprintf("split-%s-to-files", resourceName),
		false,
		"Should each resource be separate TF file",
	)
}
