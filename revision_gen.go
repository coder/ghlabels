// +build ignore

package main

import (
	"context"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"
)

func main() {
	log.SetFlags(0)

	revision := revision()
	writeRevision(revision)
}

func revision() string {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	git := exec.CommandContext(ctx, "git", "describe", "--dirty", "--always")
	revisionBytes, err := git.Output()
	if err != nil {
		log.Fatalf("failed to run `git describe`: %v", err)
	}

	revision := string(revisionBytes)
	return strings.TrimSpace(revision)
}

func writeRevision(revision string) {
	file := revisionFile(revision)

	formattedFile, err := format.Source(file)
	if err != nil {
		log.Fatalf("failed to format revision file %q: %v", string(file), err)
	}

	err = ioutil.WriteFile("revision.go", formattedFile, 0644)
	if err != nil {
		log.Fatalf("failed to write revision.go: %v", err)
	}
}

func revisionFile(revision string) []byte {
	tmpl := `package main

//go:generate go run revision_gen.go
var revision = %q
`

	return []byte(fmt.Sprintf(tmpl, revision))
}
