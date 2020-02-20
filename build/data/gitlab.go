package data

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strings"
)

type (
	releaseInfo struct {
		TagName     string `json:"tag_name"`
		Description string `json:"description"`
	}

	tagInfo struct {
		Name string `json:"name"`
	}
)

//GetGitlabReleases returns all Sia release tags matching the format v0.0.0 or v0.0.0-rc0
func GetGitlabReleases() (tags []string, latest string, err error) {
	var releaseMeta []tagInfo

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
		if versionCmp(release.Name, "v1.4.1") == -1 {
			continue
		}

		if !strings.Contains(release.Name, "-") && versionCmp(release.Name, latest) == 1 {
			latest = release.Name
		}

		tags = append(tags, release.Name)
	}

	sort.Slice(tags, func(i, j int) bool {
		if versionCmp(tags[i], tags[j]) == -1 {
			return true
		}

		return false
	})

	return
}
