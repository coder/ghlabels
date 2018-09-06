package main

import (
	"github.com/google/go-github/github"
)

func makeLabels(ghlabels []*github.Label) []*label {
	labels := make([]*label, 0, len(ghlabels))
	for _, gl := range ghlabels {
		labels = append(labels, fromGithubLabel(gl))
	}
	return labels
}

type label struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
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
