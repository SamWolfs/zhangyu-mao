package config

import (
	"github.com/spf13/viper"
)

type GitHubAppConfig struct {
	AppId          int
	PrivateKeyPath string
	WebhookSecret  string
}

var gitHubAppConfig *GitHubAppConfig

func GetGitHubAppConig() GitHubAppConfig {
	if gitHubAppConfig == nil {
		gitHubAppConfig = &GitHubAppConfig{
			AppId:          viper.GetInt("APP_ID"),
			PrivateKeyPath: viper.GetString("PRIVATE_KEY_PATH"),
			WebhookSecret:  viper.GetString("WEBHOOK_SECRET"),
		}
	}
	return *gitHubAppConfig
}
