package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type GitHubAppConfig struct {
	AppId         int
	PrivateKey    []byte
	WebhookSecret string
}

var gitHubAppConfig *GitHubAppConfig

func GetGitHubAppConig() GitHubAppConfig {
	if gitHubAppConfig == nil {
		privateKeyPath := viper.GetString("PRIVATE_KEY_PATH")
		privateKey, err := os.ReadFile(privateKeyPath)

		if err != nil {
			error := fmt.Sprintf("Unable to read GitHub Private Key (%s): %s", privateKeyPath, err)
			panic(error)
		}

		gitHubAppConfig = &GitHubAppConfig{
			AppId:         viper.GetInt("APP_ID"),
			PrivateKey:    privateKey,
			WebhookSecret: viper.GetString("WEBHOOK_SECRET"),
		}
	}
	return *gitHubAppConfig
}
