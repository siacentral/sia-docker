package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
)

type (
	dockerTagInfo struct {
		Count   uint64    `json:"count"`
		Next    *string   `json:"next"`
		Results []tagInfo `json:"results"`
	}
)

// GetDockerTags returns all release tags matching the format v0.0.0-rc0
func GetDockerTags(repo string) (tags []string, err error) {
	getTags := func(url string) (*string, error) {
		var releaseInfo dockerTagInfo

		req, err := http.NewRequest("GET", url, nil)

		if err != nil {
			return nil, err
		}

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Encoding", "deflate, gzip")

		resp, err := client.Do(req)

		if err != nil {
			return nil, err
		}

		defer drainAndClose(resp.Body)

		if err = json.NewDecoder(resp.Body).Decode(&releaseInfo); err != nil {
			return nil, err
		}

		if len(releaseInfo.Results) == 0 {
			err = errors.New("no releases")
			return nil, err
		}

		for _, tag := range releaseInfo.Results {
			if !dockerReg.MatchString(tag.Name) {
				continue
			}

			tags = append(tags, getVersion(tag.Name))
		}

		return releaseInfo.Next, nil
	}

	url := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags", repo)

	for {
		next, err := getTags(url)

		if err != nil {
			return []string{}, err
		}

		if next == nil {
			break
		}

		url = *next
	}

	sort.Slice(tags, func(i, j int) bool {
		if VersionCmp(tags[i], tags[j]) == -1 {
			return true
		}

		return false
	})

	return
}
