package gh

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"time"
)

func Repos(ctx context.Context, gc *github.Client, org string) ([]*github.Repository, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	opts := &github.RepositoryListByOrgOptions{}
	var repos []*github.Repository
	for {
		repos2, resp, err := gc.Repositories.ListByOrg(ctx, org, opts)
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
	if err != nil {
		return fmt.Errorf("failed to delete label %q on repo %v/%v: %v", label, owner, repo, err)
	}
	return nil
}

func EditLabel(ctx context.Context, gc *github.Client, owner, repo string, labelName string,
	label *github.Label) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	_, _, err := gc.Issues.EditLabel(ctx, owner, repo, labelName, label)
	if err != nil {
		return fmt.Errorf("failed to edit label %q on repo %v/%v: %v", label.GetName(), owner, repo, err)
	}
	return nil
}

func CreateLabel(ctx context.Context, gc *github.Client, owner, repo string, label *github.Label) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	_, _, err := gc.Issues.CreateLabel(ctx, owner, repo, label)
	if err != nil {
		return fmt.Errorf("failed to create label %q on repo %v/%v: %v", label.GetName(), owner, repo, err)
	}
	return nil
}

func Labels(ctx context.Context, gc *github.Client, owner, repo string) ([]*github.Label, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	opts := &github.ListOptions{}
	var ls []*github.Label
	for {
		labels2, resp, err := gc.Issues.ListLabels(ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list labels for repo %v/%v: %v", owner, repo, err)
		}
		ls = append(ls, labels2...)
		if resp.NextPage == 0 {
			return ls, nil
		}
		opts.Page = resp.NextPage
	}
}
