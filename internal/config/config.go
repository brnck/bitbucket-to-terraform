package config

type Config struct {
	BitbucketUsername  string
	BitbucketPassword  string
	BitbucketWorkspace string
	DryRun             bool
	LogLevel           uint
	Projects           ResourceFetchConfig
	Repositories       ResourceFetchConfig
}

type ResourceFetchConfig struct {
	Fetch        bool
	Path         string
	SplitToFiles bool
}
