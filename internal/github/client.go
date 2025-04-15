package github

import (
	"net/http"
	"zhangyumao/config"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v71/github"
)

type Client = github.Client
type Installation = github.Installation

// NewAppsClient returns a new github.Client that performs app authentication for
// the GitHub App using its appId and private key.
// The client can only be used for performing app-level operations that are not associated
// with a specific installation.
//
// Authenticating as a GitHub App lets you do a couple of things:
//  * You can retrieve high-level management information about your GitHub App.
//  * You can request access tokens for an installation of the app.
//
// Tips for determining the arguments for this function:
//  * the integration ID is listed as "ID" in the "About" section of the app's page
//  * the key bytes must be a PEM-encoded PKCS1 or PKCS8 private key for the application
func NewAppsClient(config config.GitHubAppConfig) (*Client, error) {
	tr := http.DefaultTransport

	itr, err := ghinstallation.NewAppsTransportKeyFromFile(tr, int64(config.AppId), config.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	return github.NewClient(&http.Client{Transport: itr}), nil
}
