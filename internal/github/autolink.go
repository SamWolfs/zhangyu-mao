package github

import (
	"context"
	"fmt"
	"slices"
	"zhangyumao/internal/errors"

	"github.com/google/go-github/v71/github"
)

type AutolinksService service
type Autolink = github.AutolinkOptions

func (s *AutolinksService) Synchronise(ctx context.Context, owner, repo string, autolinks *[]Autolink) (*[]Autolink, *errors.ErrorList) {
	var errs errors.ErrorList

	if autolinks == nil {
		return nil, nil
	}

	currentAutolinks, _, err := s.Client.Repositories.ListAutolinks(ctx, owner, repo, nil)
	if err != nil {
		err = fmt.Errorf("error while retrieving existing autolinks: %v", err)
		errs.Push(err)
		return nil, &errs
	}

	var configuredPrefixes []string
	for _, autolink := range currentAutolinks {
		configuredPrefixes = append(configuredPrefixes, *autolink.KeyPrefix)
	}

	var newAutolinks []Autolink
	for _, autolink := range *autolinks {
		if !slices.Contains(configuredPrefixes, *autolink.KeyPrefix) {
			newAutolinks = append(newAutolinks, autolink)
			_, _, err := s.Client.Repositories.AddAutolink(ctx, owner, repo, &autolink)
			if err != nil {
				err = fmt.Errorf("error while adding autlink: %v", err)
				errs.Push(err)
			}
		}
	}

	return &newAutolinks, &errs
}
