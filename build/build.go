package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/siacentral/sia-docker/build/data"
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

func handleRelease(commit string, tags ...string) (successful []string, err error) {
	log.Printf("Building %s from %s", strings.Join(tags, ", "), commit)

	builtTags := []string{}
	buildArgs := []string{"buildx",
		"build",
		"--no-cache",
		"--build-arg",
		fmt.Sprintf("SIA_VERSION=%s", commit),
		"--platform",
		"linux/amd64,linux/arm64",
		"--push"}

	for _, tag := range tags {
		tag = fmt.Sprintf("%s:%s", dockerHubRepo, tag)
		buildArgs = append(buildArgs, "-t", tag)
		builtTags = append(builtTags, tag)
	}

	buildArgs = append(buildArgs, ".")
	err = runCommand(dockerPath, buildArgs...)
	if err != nil {
		return
	}

	successful = append(successful, builtTags...)

	return
}

func buildDocker() {
	releases, latest, lastRC, err := data.GetGitlabReleases()

	if err != nil {
		log.Fatalln(err)
	}

	built, err := data.GetDockerTags(dockerHubRepo)

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
		log.Println(tag)

		// skip release if it's already found on docker and we're not overwriting
		if !overwrite && builtTags[tag] {
			continue
		}

		tags := []string{tag}

		if tag == latest {
			tags = append(tags, "latest")

			// if the lastRC is before the latest official release update the beta tag
			if data.VersionCmp(lastRC, latest) == -1 {
				tags = append(tags, "beta")
			}
		}

		// tags are normalized without the leading "v", so we need to add it for the commit id
		pushed, err := handleRelease("v"+tag, tags...)
		if err != nil {
			log.Fatalln(err)
		}

		successfulTags = append(pushed, successfulTags...)
	}

	// rebuild the latest RC because the commit id has changed in the past
	// also tags the release with "beta" since it is the latest RC
	if data.VersionCmp(lastRC, latest) == 1 {
		pushed, err := handleRelease("v"+lastRC, lastRC, "beta")
		if err != nil {
			log.Fatalln(err)
		}

		successfulTags = append(pushed, successfulTags...)
	}

	//build the unstable master branch
	pushed, err := handleRelease("master", "unstable")
	if err != nil {
		log.Fatalln(err)
	}

	successfulTags = append(successfulTags, pushed...)

	//log pushed tags
	log.Println("Successfully built and pushed:", strings.Join(successfulTags, ", "))
}

func main() {
	flag.StringVar(&dockerHubRepo, "docker-hub-repo", "", "the docker hub repository to push to")
	flag.StringVar(&dockerPath, "docker-path", "/usr/bin/docker", "the path to docker")
	flag.BoolVar(&overwrite, "overwrite", false, "overwrite existing tags with new builds")
	flag.Parse()

	if len(dockerHubRepo) == 0 {
		log.Fatalln("--docker-hub-repo is required")
	}

	buildDocker()
}
