package github

import (
	"net/http"
	"zhangyumao/config"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v71/github"
)

// TODO: Consider turning this into a ClientCreator/ClientFactory interface
// which should(?) allow mocking of the github Clients in e.g. GenServers
type Client struct {
	github.Client

	Autolinks *AutolinksService
}
type Installation = github.Installation
type ClientCreator interface {
	NewAppClient() (*Client, error)
	NewInstallationClient(int64) (*Client, error)
}

type PushEvent = github.PushEvent

func NewClientCreator(c config.GitHubAppConfig) clientCreator {
	return clientCreator{
		config: c,
	}
}

type clientCreator struct {
	config config.GitHubAppConfig
}

var _ ClientCreator = &clientCreator{}

// Descriptions partially copied from
// [palantor/go-githubapp](https://github.com/palantir/go-githubapp/blob/develop/githubapp/client_creator.go)

// NewAppsClient returns a new github.Client that performs app authentication for
// the GitHub App using its appId and private key.
// The client can only be used for performing app-level operations that are not associated
// with a specific installation.
//
// Authenticating as a GitHub App lets you do a couple of things:
//   - You can retrieve high-level management information about your GitHub App.
//   - You can request access tokens for an installation of the app.
//
// Tips for determining the arguments for this function:
//   - the integration ID is listed as "ID" in the "About" section of the app's page
//   - the key bytes must be a PEM-encoded PKCS1 or PKCS8 private key for the application
func (c *clientCreator) NewAppClient() (*Client, error) {
	tr := http.DefaultTransport

	itr, err := ghinstallation.NewAppsTransportKeyFromFile(tr, c.config.AppId, c.config.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	client := &Client{Client: *github.NewClient(&http.Client{Transport: itr})}
	client.initialize()

	return client, nil
}

// NewInstallationClient returns a new github.Client that performs app authentication for
// the GitHub App using its appId, private key, and the given installation ID.
// The client can be used to perform all operations that the GitHub App is configured for.
//
// Tips for determining the arguments for this function:
//   - the integration ID is listed as "ID" in the "About" section of the app's page
//   - the installation ID is the ID that is shown in the URL of https://{githubURL}/settings/installations/{#}
//     (navigate to the "installations" page without the # and go to the app's page to see the number)
//   - the key bytes must be a PEM-encoded PKCS1 or PKCS8 private key for the application
func (c *clientCreator) NewInstallationClient(installationID int64) (*Client, error) {
	tr := http.DefaultTransport

	itr, err := ghinstallation.NewKeyFromFile(tr, c.config.AppId, installationID, c.config.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	client := &Client{Client: *github.NewClient(&http.Client{Transport: itr})}
	client.initialize()

	return client, nil
}

func (c *Client) initialize() {
	c.Autolinks = &AutolinksService{Client: &c.Client}
}
