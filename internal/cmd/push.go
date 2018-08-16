package cmd

import (
	"github.com/spf13/cobra"
	"encoding/json"
	"os"
	"log"
	"github.com/nhooyr/labels/internal/gh"
	"strings"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push labels to a remote github repository",
	Long: `push lets you push a set of labels to a remote github repository.
It will read the set of labels from stdin in the same format as
what is output from 'labels pull'. However, it will allow you to set
a from field for each label indicating to replace the from label
with the data specified by the new one. However, if the from label does
not exist, push will just ensure that the new label exists with the correct
specification. If both exist, then an error will be thrown.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var org string
		var repos []string

		if strings.ContainsRune(args[0], '/') {
			repoPath := args[0]
			var repo string
			var err error
			org, repo, err = splitRepoPath(repoPath)
			if err != nil {
				return err
			}
			repos = []string{repo}
		} else {
			org = args[0]
			ghrepos, err := gh.Repos(ctx, gc, org)
			if err != nil {
				log.Fatalf("failed to list repos for org: %v", err)
			}

			for _, gr := range ghrepos {
				repos = append(repos, gr.GetName())
			}
		}

		var labels labels
		d := json.NewDecoder(os.Stdin)
		d.DisallowUnknownFields()
		err := d.Decode(&labels)
		if err != nil {
			log.Fatalf("failed to decode input labels: %v", err)
		}

		fromMap, err := labels.makeFromMap()
		if err != nil {
			log.Fatalf("failed to make from map: %v", err)
		}

		for _, repo := range repos {
			ghlabels, err := gh.Labels(ctx, gc, org, repo)
			if err != nil {
				log.Fatal(err)
			}

			labelsMap, err := labels.makeMap()
			if err != nil {
				log.Fatalf("failed to make labels map: %v", err)
			}

			for _, gl := range ghlabels {
				var l *label

				l1, ok1 := labelsMap[gl.GetName()]
				l2, ok2 := fromMap[gl.GetName()]

				if ok1 {
					if ok2 {
						log.Fatalf("label %v cannot be from label %v in config", l2.Name, gl.GetName())
					}
					l = l1
				} else if ok2 {
					l = l2
				} else {
					if keepDefaults {
						continue
					}
					if gl.GetDefault() {
						err = gh.DeleteLabel(ctx, gc, org, repo, gl.GetName())
						if err != nil {
							log.Fatalf("failed to delete default label: %v", err)
						}
					}
					continue
				}

				err = gh.EditLabel(ctx, gc, org, repo, gl.GetName(), l.github())
				if err != nil {
					log.Fatalf("failed to edit label: %v", err)
				}

				delete(labelsMap, l.Name)
			}

			for _, l := range labelsMap {
				err = gh.CreateLabel(ctx, gc, org, repo, l.github())
				if err != nil {
					log.Fatalf("failed to create label: %v", err)
				}

				delete(labelsMap, l.Name)
			}

			if len(labelsMap) > 0 {
				log.Fatalf("could not ensure all labels on repo %v/%v for unknown reason; left: %+v", org, repo,
					labelsMap)
			}
		}

		return nil
	},
}

var keepDefaults bool

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.Flags().BoolVar(&keepDefaults, "keep-defaults", false, "Keep default labels when pushing")
}
