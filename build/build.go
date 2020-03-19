package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/siacentral/docker-sia/build/data"
)

var (
	dockerPath    string
	dockerHubRepo string
	overwrite     bool
)

func handleOutput(out io.Reader) {
	in := bufio.NewScanner(out)

	for in.Scan() {
		text := in.Text()

		log.Println(text)
	}
}

func runCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	go handleOutput(stdOut)

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go handleOutput(stdErr)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func handleRelease(tag, commit string, latest bool) (successful []string, err error) {
	log.Printf("Building %s from %s", tag, commit)

	dockerTag := fmt.Sprintf("%s:%s", dockerHubRepo, tag)
	buildArgs := []string{"build",
		"--build-arg",
		fmt.Sprintf("SIA_VERSION=%s", commit),
		"-t", dockerTag}

	// if this is the latest full release tag it with latest too
	if latest {
		buildArgs = append(buildArgs, "-t", fmt.Sprintf("%s:latest", dockerHubRepo))
	}

	buildArgs = append(buildArgs, ".")
	err = runCommand(dockerPath, buildArgs...)
	if err != nil {
		return
	}

	err = runCommand(dockerPath, "push", dockerTag)
	if err != nil {
		return nil, err
	}

	successful = append(successful, tag)

	// if this is the latest tag push it to docker hub too
	if latest {
		err = runCommand(dockerPath, "push", fmt.Sprintf("%s:latest", dockerHubRepo))
		if err != nil {
			return
		}
		successful = append(successful, fmt.Sprintf("latest (%s)", tag))
	}

	return
}

func main() {
	flag.StringVar(&dockerHubRepo, "docker-hub-repo", "", "the docker hub repository to push to")
	flag.StringVar(&dockerPath, "docker-path", "/usr/bin/docker", "the path to docker")
	flag.BoolVar(&overwrite, "overwrite", false, "overwrite existing tags with new builds")
	flag.Parse()

	if len(dockerHubRepo) == 0 {
		log.Fatalln("--docker-hub-repo is required")
	}

	releases, latest, err := data.GetGitlabReleases()

	if err != nil {
		log.Fatalln(err)
	}

	built, err := data.GetDockerTags()

	if err != nil {
		log.Fatalln(err)
	}

	builtTags := make(map[string]bool)

	for _, tag := range built {
		builtTags[tag] = true
	}

	successfulTags := []string{}

	// loop through all found releases
	for _, tag := range releases {
		// skip release if it's already found on docker and we're not overwriting
		if !overwrite && builtTags[tag] {
			continue
		}

		// tags are normalized without the leading "v", so we need to add it for the commit id
		pushed, err := handleRelease(tag, "v"+tag, tag == latest)
		if err != nil {
			log.Fatalln(err)
		}

		successfulTags = append(pushed, successfulTags...)
	}

	//build the unstable master branch
	pushed, err := handleRelease("unstable", "master", false)
	if err != nil {
		log.Fatalln(err)
	}

	successfulTags = append(successfulTags, pushed...)

	//log pushed tags
	log.Println("Successfully built and pushed:", strings.Join(successfulTags, ", "))
}
