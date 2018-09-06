// +build ignore

package ghlabels

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

	version := version()
	writeVersion(version)
}

func version() string {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	git := exec.CommandContext(ctx, "git", "describe", "--dirty", "--always")
	versionBytes, err := git.Output()
	if err != nil {
		log.Fatalf("failed to run `git describe`: %v", err)
	}

	version := string(versionBytes)
	return strings.TrimSpace(version)
}

func writeVersion(version string) {
	file := versionFile(version)

	formattedFile, err := format.Source(file)
	if err != nil {
		log.Fatalf("failed to format version file %q: %v", string(file), err)
	}

	err = ioutil.WriteFile("version.go", formattedFile, 0644)
	if err != nil {
		log.Fatalf("failed to write version.go: %v", err)
	}
}

func versionFile(version string) []byte {
	tmpl := `package main

//go:generate go run version_gen.go
var version = %q
`

	return []byte(fmt.Sprintf(tmpl, version))
}
