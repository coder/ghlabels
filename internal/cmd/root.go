package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/google/go-github/github"
	hubgh "github.com/github/hub/github"
	"context"
	"golang.org/x/oauth2"
	"net/url"
	"log"
	"errors"
)

var gc *github.Client
var ctx = context.Background()

var rootCmd = &cobra.Command{
	Use:     "labels",
	Short:   "Synchronize github labels",
	Long:    `labels lets you synchronize github labels sanely.`,
	Version: version,
	PersistentPreRunE: func(_ *cobra.Command, args []string) error {
		log.SetFlags(0)

		if len(args) < 1 {
			return errors.New("need at least one arg")
		}

		config := hubgh.CurrentConfig()
		h, err := config.DefaultHost()
		if err != nil {
			log.Fatalf("failed to get host from hub: %v", err)
		}

		ctx := context.Background()
		hc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: h.AccessToken,
		}))
		gc = github.NewClient(hc)

		gc.BaseURL = &url.URL{}
		gc.BaseURL.Scheme = h.Protocol

		if h.Host == "github.com" {
			gc.BaseURL.Host = "api.github.com"
			gc.BaseURL.Path = "/"
		} else {
			gc.BaseURL.Host = h.Host
			gc.BaseURL.Path = "/api/v3/"
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
