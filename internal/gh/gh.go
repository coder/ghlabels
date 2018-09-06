// gh is a wrapper package around github.com/google/go-github for ghlabels.
package gh

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

func UserRepos(ctx context.Context, gc *github.Client) ([]*github.Repository, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	opts := &github.RepositoryListOptions{
		Type: "owner",
	}
	var repos []*github.Repository
	for {
		repos2, resp, err := gc.Repositories.List(ctx, "", opts)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repos2...)
		if resp.NextPage == 0 {
			return repos, nil
		}
		opts.Page = resp.NextPage
	}
}

func OrgRepos(ctx context.Context, gc *github.Client, org string) ([]*github.Repository, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	opts := github.RepositoryListByOrgOptions{}
	var repos []*github.Repository
	for {
		repos2, resp, err := gc.Repositories.ListByOrg(ctx, org, &opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list repos for %v: %v", org, err)
		}
		repos = append(repos, repos2...)
		if resp.NextPage == 0 {
			return repos, nil
		}
		opts.Page = resp.NextPage
	}
}

func DeleteLabel(ctx context.Context, gc *github.Client, owner, repo, label string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	_, err := gc.Issues.DeleteLabel(ctx, owner, repo, label)
	return err
}

func EditLabel(ctx context.Context, gc *github.Client, owner, repo, labelName string,
	label *github.Label) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	_, _, err := gc.Issues.EditLabel(ctx, owner, repo, labelName, label)
	return err
}

var ErrAlreadyExists = errors.New("label already exists")

func CreateLabel(ctx context.Context, gc *github.Client, owner, repo string, label *github.Label) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	_, _, err := gc.Issues.CreateLabel(ctx, owner, repo, label)
	if err != nil {
		er, ok := err.(*github.ErrorResponse)
		if ok && len(er.Errors) > 0 {
			if er.Errors[0].Code == "already_exists" {
				return ErrAlreadyExists
			}
		}
	}
	return err
}

func Labels(ctx context.Context, gc *github.Client, owner, repo string) ([]*github.Label, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	opts := &github.ListOptions{}
	var ls []*github.Label
	for {
		labels2, resp, err := gc.Issues.ListLabels(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}
		ls = append(ls, labels2...)
		if resp.NextPage == 0 {
			return ls, nil
		}
		opts.Page = resp.NextPage
	}
}
