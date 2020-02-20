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

func main() {
	var dockerPath, dockerHubRepo string
	var overwrite bool

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
		builtTags[fmt.Sprintf("v%s", tag)] = true
	}

	successfulTags := []string{}

	// loop through all found releases
	for _, tag := range releases {
		// skip release if it's already found on docker and we're not overwriting
		if !overwrite && builtTags[tag] {
			continue
		}

		log.Println("Building ", tag)
		dockerTag := fmt.Sprintf("%s:%s", dockerHubRepo, tag[1:])
		buildArgs := []string{"build",
			"--build-arg",
			fmt.Sprintf("SIA_VERSION=%s", tag),
			"-t", dockerTag}

		// if this is the latest full release tag it with latest too
		if tag == latest {
			buildArgs = append(buildArgs, "-t", fmt.Sprintf("%s:latest", dockerHubRepo))
		}

		buildArgs = append(buildArgs, ".")
		err := runCommand(dockerPath, buildArgs...)
		if err != nil {
			log.Fatalln(err)
		}

		err = runCommand(dockerPath, "push", dockerTag)
		if err != nil {
			log.Fatalln(err)
		}

		successfulTags = append(successfulTags, tag)

		// if this is the latest tag push it to docker hub too
		if tag == latest {
			err = runCommand(dockerPath, "push", fmt.Sprintf("%s:latest", dockerHubRepo))
			if err != nil {
				log.Fatalln(err)
			}
			successfulTags = append(successfulTags, fmt.Sprintf("latest (%s)", tag))
		}

	}

	log.Println("Building unstable")

	//build the unstable master branch
	err = runCommand(dockerPath, "build",
		"--build-arg",
		"SIA_VERSION=master",
		"-t",
		fmt.Sprintf("%s:unstable", dockerHubRepo),
		".")

	if err != nil {
		log.Fatalln(err)
	}

	err = runCommand(dockerPath, "push", fmt.Sprintf("%s:unstable", dockerHubRepo))
	if err != nil {
		log.Fatalln(err)
	}

	successfulTags = append(successfulTags, "unstable")

	//log pushed tags
	if len(successfulTags) > 0 {
		log.Println("Successfully built and pushed:", strings.Join(successfulTags, ", "))
		return
	}

	log.Fatalln("No new releases built")
}
