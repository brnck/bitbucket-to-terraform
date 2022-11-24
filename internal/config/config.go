package config

// Config application configuration
type Config struct {
	BitbucketUsername        string
	BitbucketPassword        string
	BitbucketWorkspace       string
	LogLevel                 int
	GenerateImportStatements bool
	ImportStatementsPath     string
	Projects                 ResourceFetchConfig
	Repositories             ResourceFetchConfig
}

// ResourceFetchConfig each resource from bitbucket can be fetched differently or not fetched at all
type ResourceFetchConfig struct {
	Fetch bool
	Path  string
}
