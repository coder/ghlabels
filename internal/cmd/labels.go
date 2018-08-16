package cmd

import (
	"strings"
	"fmt"
	"github.com/google/go-github/github"
)

func splitRepoPath(repoPath string) (owner, repo string, err error) {
	s := strings.Split(repoPath, "/")
	if len(s) > 2 {
		return "", "", fmt.Errorf("too many slashes in repo path %q", repoPath)
	}

	return s[0], s[1], nil
}

type labels []*label

func (ls labels) makeMap() (map[string]*label, error) {
	labelsMap := make(map[string]*label, len(ls))
	for _, l := range ls {
		l2, ok := labelsMap[l.Name]
		if ok {
			return nil, fmt.Errorf("two labels with name %v", l2.Name)
		}
		labelsMap[l.Name] = l
	}
	return labelsMap, nil
}

func (ls labels) makeFromMap() (map[string]*label, error) {
	fromMap := make(map[string]*label)
	for _, l := range ls {
		if l.From != "" {
			l2, ok := fromMap[l.From]
			if ok {
				return nil, fmt.Errorf("labels %v and %v are both from %v?", l2.Name, l.Name, l.From)
			}
			fromMap[l.From] = l
		}
	}
	return fromMap, nil
}

func makeLabels(ghlabels []*github.Label) labels {
	labels := make(labels, 0, len(ghlabels))
	for _, gl := range ghlabels {
		labels = append(labels, fromGithubLabel(gl))
	}
	return labels
}

type label struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	From        string `json:"from,omitempty"`
}

func fromGithubLabel(gl *github.Label) *label {
	l := &label{
		Name:        gl.GetName(),
		Description: gl.GetDescription(),
		Color:       gl.GetColor(),
	}
	return l
}

func (l *label) github() *github.Label {
	return &github.Label{
		Name:        github.String(l.Name),
		Description: github.String(l.Description),
		Color:       github.String(l.Color),
	}
}
