package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type (
	// GitlabRelease info about the latest releases
	GitlabRelease struct {
		Name   string `json:"name"`
		Target string `json:"target"`
	}

	commitInfo struct {
		ID string `json:"id"`
	}
)

// GetLastCommit returns the last commit of the target ref
func GetLastCommit(ref string) (string, error) {
	var commitMeta []commitInfo

	req, err := http.NewRequest("GET", fmt.Sprintf("https://gitlab.com/api/v4/projects/7508674/repository/commits?ref=%s", url.QueryEscape(ref)), nil)

	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Encoding", "deflate, gzip")

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer drainAndClose(resp.Body)

	if err = json.NewDecoder(resp.Body).Decode(&commitMeta); err != nil {
		return "", err
	}

	if len(commitMeta) == 0 {
		return "", errors.New("no commits")
	}

	return commitMeta[0].ID, nil
}

//GetGitlabReleases returns all Sia release tags matching the format v0.0.0 or v0.0.0-rc0
func GetGitlabReleases() (tags []GitlabRelease, latest GitlabRelease, lastRC GitlabRelease, err error) {
	var releaseMeta []GitlabRelease

	req, err := http.NewRequest("GET", "https://gitlab.com/api/v4/projects/7508674/repository/tags?order_by=name", nil)

	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Encoding", "deflate, gzip")

	resp, err := client.Do(req)

	if err != nil {
		return
	}

	defer drainAndClose(resp.Body)

	if err = json.NewDecoder(resp.Body).Decode(&releaseMeta); err != nil {
		return
	}

	if len(releaseMeta) == 0 {
		err = errors.New("no releases")
		return
	}

	for _, release := range releaseMeta {
		if !gitlabRegex.MatchString(release.Name) {
			continue
		}

		// versions below v1.4.1 won't build properly with the Dockerfile
		if VersionCmp(release.Name, "1.4.1") == -1 {
			continue
		}

		tag := getVersion(release.Name)

		// check for latest RC release and latest official release
		if strings.Contains(release.Name, "-") && VersionCmp(release.Name, lastRC.Name) == 1 {
			lastRC.Name = tag
			lastRC.Target = release.Target
		} else if !strings.Contains(release.Name, "-") && VersionCmp(release.Name, latest.Name) == 1 {
			latest.Name = tag
			latest.Target = release.Target
		}

		tags = append(tags, GitlabRelease{
			Name:   tag,
			Target: release.Target,
		})
	}

	sort.Slice(tags, func(i, j int) bool {
		if VersionCmp(tags[i].Name, tags[j].Name) == -1 {
			return true
		}

		return false
	})

	return
}
