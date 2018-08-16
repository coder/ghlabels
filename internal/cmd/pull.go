package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/nhooyr/labels/internal/gh"
	"encoding/json"
	"os"
	"log"
)

var pullCmd = &cobra.Command{
	Use:   "pull [repo]/[path]",
	Short: "Pull labels from a remote repository",
	Long: `Pull will pull labels from a remote repository and output them
to stdout as a JSON map from the label name to a label structure.
You can use this structure later with the 'labels push' command to push
labels to a remote repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repoPath := args[0]

		owner, repo, err := splitRepoPath(repoPath)
		if err != nil {
			return fmt.Errorf("failed to split repo path %v: %v", repoPath, err)
		}

		ghlabels, err := gh.Labels(ctx, gc, owner, repo)
		if err != nil {
			log.Fatalf("failed to pull %v: %v", repoPath, err)
		}

		labels := makeLabels(ghlabels)

		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "    ")
		err = e.Encode(labels)
		if err != nil {
			log.Fatalf("failed to encode labels to stdout: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
