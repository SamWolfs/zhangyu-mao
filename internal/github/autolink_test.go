package github

import (
	"context"
	"testing"

	"github.com/google/go-github/v71/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
)

func TestAutolinksService_Synchronise(t *testing.T) {
	client := testAutolinksService_MockClient()
	service := AutolinksService{Client: client}
	ctx := context.Background()
	owner, repo := "owner", "repo"

	autolinks := []Autolink{
		{
			KeyPrefix:      github.Ptr("PREFIX2-"),
			URLTemplate:    github.Ptr("https://planning-tool.com/PREFIX2-<num>"),
			IsAlphanumeric: github.Ptr(true),
		},
		{
			KeyPrefix:      github.Ptr("PREFIX3-"),
			URLTemplate:    github.Ptr("https://planning-tool.com/PREFIX3-<num>"),
			IsAlphanumeric: github.Ptr(true),
		},
	}

	got, _ := service.Synchronise(ctx, owner, repo, &autolinks)
	want := &[]Autolink{autolinks[1]}

	assert.Equal(t, want, got)
}

func testAutolinksService_MockClient() *github.Client {
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposAutolinksByOwnerByRepo,
			[]github.Autolink{
				{
					ID:             github.Ptr(int64(1)),
					KeyPrefix:      github.Ptr("PREFIX1-"),
					URLTemplate:    github.Ptr("https://planning-tool.com/PREFIX1-<num>"),
					IsAlphanumeric: github.Ptr(true),
				},
				{
					ID:             github.Ptr(int64(2)),
					KeyPrefix:      github.Ptr("PREFIX2-"),
					URLTemplate:    github.Ptr("https://planning-tool.com/PREFIX2-<num>"),
					IsAlphanumeric: github.Ptr(true),
				},
			}))

	return github.NewClient(mockedHTTPClient)
}
