package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/siacentral/sia-docker/build/data"
)

var (
	dockerPath    string
	dockerHubRepo string
	built         = make(map[string]string)
)

func handleOutput(f io.Writer, out io.Reader) {
	in := bufio.NewScanner(out)

	for in.Scan() {
		text := in.Text()

		_, err := f.Write([]byte(text + "\n"))
		if err != nil {
			log.Println("err writing log:", err)
		}
	}
}

func runCommand(logPath string, command string, args ...string) error {
	f, err := os.Create(logPath)
	if err != nil {
		return err
	}

	defer f.Close()

	cmd := exec.Command(command, args...)

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	go handleOutput(f, stdOut)

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go handleOutput(f, stdErr)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func handleRelease(commit string, tags ...string) (successful []string, err error) {
	log.Printf("Building %s from %s", strings.Join(tags, ", "), commit)

	start := time.Now()
	builtTags := []string{}
	buildArgs := []string{"buildx",
		"build",
		"--no-cache",
		"--build-arg",
		fmt.Sprintf("SIA_VERSION=%s", commit),
		"--build-arg",
		fmt.Sprintf("RC=\"%s\"", data.GetRC(commit)),
		"--platform",
		"linux/amd64,linux/arm64,linux/arm/v6,linux/arm/v7",
		"--push"}

	for _, tag := range tags {
		tag = fmt.Sprintf("%s:%s", dockerHubRepo, tag)
		buildArgs = append(buildArgs, "-t", tag)
		builtTags = append(builtTags, tag)
	}

	buildArgs = append(buildArgs, ".")
	err = runCommand(fmt.Sprintf("%s.log", commit), dockerPath, buildArgs...)
	if err != nil {
		return
	}

	successful = append(successful, builtTags...)

	log.Printf(" Finished in %s", time.Since(start))

	return
}

func buildDocker() error {
	releases, latest, lastRC, err := data.GetGitlabReleases()

	if err != nil {
		return fmt.Errorf("error getting releases: %s", err)
	}

	successfulTags := []string{}

	// loop through all found releases
	for _, release := range releases {
		// skip release only if it matches the commit id. Since these aren't persisted this will now always build all releases on the first run.
		if built[release.Name] == release.Target {
			continue
		}

		tags := []string{release.Name}

		if release.Name == latest.Name {
			tags = append(tags, "latest")

			// if the lastRC is before the latest official release update the beta tag
			if data.VersionCmp(lastRC.Name, latest.Name) == -1 {
				tags = append(tags, "beta")
			}
		}

		// tags are normalized without the leading "v", so we need to add it for the commit id
		pushed, err := handleRelease("v"+release.Name, tags...)
		if err != nil {
			return fmt.Errorf("error building %s: %s", release.Name, err)
		}

		successfulTags = append(pushed, successfulTags...)
		built[release.Name] = release.Target
	}

	// check that we do not need to rebuild the last RC
	if data.VersionCmp(lastRC.Name, latest.Name) == 1 && built[lastRC.Name] != lastRC.Target {
		pushed, err := handleRelease("v"+lastRC.Name, lastRC.Name, "beta")
		if err != nil {
			return fmt.Errorf("error building rc: %s", err)
		}

		successfulTags = append(pushed, successfulTags...)
		built[lastRC.Name] = lastRC.Target
	}

	masterID, err := data.GetLastCommit("master")
	if err != nil {
		return fmt.Errorf("error building rc: %s", err)
	}

	//build the unstable master branch if the commit id has changed
	if built["master"] != masterID {
		pushed, err := handleRelease("master", "unstable")
		if err != nil {
			return fmt.Errorf("error building unstable: %s", err)
		}

		successfulTags = append(successfulTags, pushed...)
		built["master"] = masterID
	}

	if len(successfulTags) == 0 {
		log.Println("no updated releases")
		return nil
	}

	log.Println("Successfully built and pushed:", strings.Join(successfulTags, ", "))
	return nil
}

func main() {
	flag.StringVar(&dockerHubRepo, "docker-hub-repo", "", "the docker hub repository to push to")
	flag.StringVar(&dockerPath, "docker-path", "/usr/bin/docker", "the path to docker")
	flag.Parse()

	if len(dockerHubRepo) == 0 {
		log.Fatalln("--docker-hub-repo is required")
	}

	for {
		if err := buildDocker(); err != nil {
			log.Println(err)
		}

		time.Sleep(time.Minute * 10)
	}
}
