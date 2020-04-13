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
	archPrefix    string
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

func handleRelease(tag, commit string) (successful []string, err error) {
	log.Printf("Building %s from %s", tag, commit)

	dockerTag := fmt.Sprintf("%s:%s", dockerHubRepo, tag)
	buildArgs := []string{"build",
		"--build-arg",
		fmt.Sprintf("SIA_VERSION=%s", commit),
		"-t", dockerTag}

	buildArgs = append(buildArgs, ".")
	err = runCommand(dockerPath, buildArgs...)
	if err != nil {
		return
	}

	err = runCommand(dockerPath, "push", dockerTag)
	if err != nil {
		return nil, err
	}

	successful = append(successful, dockerTag)

	return
}

func handleManifest(release string, tags []string) (err error) {
	releaseTag := fmt.Sprintf("%s:%s", dockerHubRepo, release)
	createArgs := []string{
		"manifest",
		"create",
		releaseTag,
	}

	log.Printf("Creating manifest for %s", releaseTag)

	for _, tag := range tags {
		dockerTag := fmt.Sprintf("%s:%s", dockerHubRepo, tag)
		createArgs = append(createArgs, dockerTag)

		log.Printf("  Adding %s", dockerTag)
	}

	err = runCommand("docker", createArgs...)
	if err != nil {
		return
	}

	for _, tag := range tags {
		dockerTag := fmt.Sprintf("%s:%s", dockerHubRepo, tag)
		parts := strings.Split(tag, "-")
		annotateArgs := []string{
			"manifest",
			"annotate",
			releaseTag,
			dockerTag,
			"--arch",
			parts[0],
		}

		if parts[0] == "arm64" {
			annotateArgs = append(annotateArgs, "--os", "arm64")
		}

		err = runCommand("docker", annotateArgs...)
		if err != nil {
			return
		}
	}

	pushArgs := []string{
		"manifest",
		"push",
		"--purge",
		releaseTag,
	}

	err = runCommand("docker", pushArgs...)

	return
}

func dockerHubTag(tag string) string {
	if len(archPrefix) == 0 {
		return tag
	}

	return fmt.Sprintf("%s-%s", archPrefix, tag)
}

func buildDocker() {
	releases, latest, err := data.GetGitlabReleases()

	if err != nil {
		log.Fatalln(err)
	}

	built, err := data.GetPrefixedDockerTags(archPrefix)

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

		dockerTag := tag

		if len(archPrefix) != 0 {
			dockerTag = fmt.Sprintf("%s-%s", archPrefix, tag)
		}

		// tags are normalized without the leading "v", so we need to add it for the commit id
		pushed, err := handleRelease(dockerTag, "v"+tag)
		if err != nil {
			log.Fatalln(err)
		}

		successfulTags = append(pushed, successfulTags...)

		if tag == latest {
			latestTag := "latest"

			if len(archPrefix) != 0 {
				latestTag = fmt.Sprintf("%s-%s", archPrefix, "latest")
			}

			pushed, err := handleRelease(latestTag, "v"+tag)
			if err != nil {
				log.Fatalln(err)
			}

			successfulTags = append(pushed, successfulTags...)
		}
	}

	unstableTag := "unstable"

	if len(archPrefix) != 0 {
		unstableTag = fmt.Sprintf("%s-%s", archPrefix, "unstable")
	}

	//build the unstable master branch
	pushed, err := handleRelease(unstableTag, "master")
	if err != nil {
		log.Fatalln(err)
	}

	successfulTags = append(successfulTags, pushed...)

	//log pushed tags
	log.Println("Successfully built and pushed:", strings.Join(successfulTags, ", "))
}

func buildManifest() {
	releaseTags, _, err := data.GetGitlabReleases()

	if err != nil {
		log.Fatalln(err)
	}

	releases := map[string][]string{
		"latest":   make([]string, 0),
		"unstable": make([]string, 0),
	}

	for _, release := range releaseTags {
		releases[release] = []string{}
	}

	dockerTags, err := data.GetDockerTags()

	if err != nil {
		log.Fatalln(err)
	}

	for _, tag := range dockerTags {
		parts := strings.Split(tag, "-")

		if len(parts) != 2 {
			continue
		}

		if _, exists := releases[parts[1]]; !exists {
			continue
		}

		releases[parts[1]] = append(releases[parts[1]], tag)
	}

	for release, dockerTags := range releases {
		if err := handleManifest(release, dockerTags); err != nil {
			log.Fatalln(err)
		}
	}
}

func main() {
	var manifestOnly bool

	flag.StringVar(&archPrefix, "arch", "", "the arch prefix to use for multi-arch support on Docker Hub")
	flag.StringVar(&dockerHubRepo, "docker-hub-repo", "", "the docker hub repository to push to")
	flag.StringVar(&dockerPath, "docker-path", "/usr/bin/docker", "the path to docker")
	flag.BoolVar(&overwrite, "overwrite", false, "overwrite existing tags with new builds")
	flag.BoolVar(&manifestOnly, "manifest", false, "build the manifest instead of the docker image")
	flag.Parse()

	if len(dockerHubRepo) == 0 {
		log.Fatalln("--docker-hub-repo is required")
	}

	if manifestOnly {
		buildManifest()
		return
	}

	buildDocker()
}
