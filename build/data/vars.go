package data

import (
	"net/http"
	"regexp"
	"time"
)

var (
	rcRegex     = regexp.MustCompile(`[0-9]+`)
	versionReg  = regexp.MustCompile(`([0-9]+.){2}[0-9]-rc[0-9]+|([0-9]+.){2}[0-9]`)
	dockerReg   = regexp.MustCompile(`^([0-9]+.){2}[0-9]-rc[0-9]+$|^([0-9]+.){2}[0-9]$`)
	gitlabRegex = regexp.MustCompile(`^v([0-9]+.){2}[0-9]-rc[0-9]+$|^v([0-9]+.){2}[0-9]$`)

	client = http.Client{
		Timeout: 10 * time.Second,
	}
)
