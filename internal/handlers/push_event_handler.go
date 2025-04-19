package handlers

import (
	"context"
	"fmt"
	"io"

	"zhangyumao/internal/github"
)

func HandlePushEvent(client *github.Client, event *github.PushEvent) error {
	owner, repo := *event.Repo.Owner.Name, *event.Repo.Name
	ctx := context.Background()
	if *event.Ref != fmt.Sprintf("refs/heads/%s", *event.Repo.DefaultBranch) {
		fmt.Println("Ignoring PushEvent on non-trunk branch.")
		return nil
	}

	r, _, err := client.Repositories.DownloadContents(ctx, owner, repo, ".github/settings.yaml", nil)

	if err != nil {
		return err
	}

	yamlSettings, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	settings, err := github.DecodeSettings(yamlSettings)
	if err != nil {
		return err
	}

	_, _, err = client.Repositories.Edit(ctx, owner, repo, settings.Repository.ToGitHub())
	if err != nil {
		return err
	}

	// TODO: print autolinks/take care of deletions
	_, errors := client.Autolinks.Synchronise(ctx, owner, repo, settings.Autolinks)
	if errors != nil {
		return errors.Err()
	}

	for _, protection := range *settings.Protections {
		_, _, err := client.Repositories.UpdateBranchProtection(ctx, owner, repo, protection.Name, &protection.ProtectionRequest)
		fmt.Printf("update branch protection: %v", err)
	}

	return nil
}
