package data

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
)

type (
	dockerTagInfo struct {
		Count   uint64    `json:"count"`
		Next    *string   `json:"next"`
		Results []tagInfo `json:"results"`
	}
)

//GetDockerTags returns all docker release tags
func GetDockerTags() (tags []string, err error) {
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
			return nil, nil
		}

		for _, tag := range releaseInfo.Results {
			tags = append(tags, tag.Name)
		}

		return releaseInfo.Next, nil
	}

	url := "https://hub.docker.com/v2/repositories/siacentral/sia/tags"

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

	return
}

//GetPrefixedDockerTags returns all release tags matching the format v0.0.0-rc0
func GetPrefixedDockerTags(prefix string) (filtered []string, err error) {
	if len(prefix) != 0 {
		prefix = prefix + "-"
	}

	tags, err := GetDockerTags()
	if err != nil {
		return
	}

	for _, tag := range tags {
		if len(prefix) != 0 && !strings.HasPrefix(tag, prefix) {
			continue
		}

		tag = strings.TrimPrefix(tag, prefix)

		if !dockerRegex.MatchString(tag) {
			continue
		}

		filtered = append(filtered, getVersion(tag))
	}

	sort.Slice(filtered, func(i, j int) bool {
		if versionCmp(tags[i], tags[j]) == -1 {
			return true
		}

		return false
	})

	return
}
