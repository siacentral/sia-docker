package data

import (
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var (
	versionNumberReg = regexp.MustCompile(`[^0-9\.]`)
	gitlabRegex      = regexp.MustCompile(`^v([0-9]+.)+[0-9]\-[\S]+|v([0-9]+.)+[0-9]$`)
	dockerRegex      = regexp.MustCompile(`^([0-9]+.)+[0-9]\-[\S]+|([0-9]+.)+[0-9]$`)
)

func versionCmp(a, b string) int {
	l := 0

	a = versionNumberReg.ReplaceAllString(a, "")
	b = versionNumberReg.ReplaceAllString(b, "")

	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	aLen := len(aParts)
	bLen := len(bParts)

	if aLen > bLen {
		l = len(aParts)
	} else {
		l = len(bParts)
	}

	for i := 0; i < l; i++ {
		var aPart, bPart int
		var err error

		if i < aLen {
			aPart, err = strconv.Atoi(aParts[i])

			if err != nil {
				return -1
			}
		}

		if i < bLen {
			bPart, err = strconv.Atoi(bParts[i])

			if err != nil {
				return 1
			}
		}

		if aPart < bPart {
			return -1
		} else if aPart > bPart {
			return 1
		}
	}

	return 0
}

func drainAndClose(rc io.ReadCloser) {
	io.Copy(ioutil.Discard, rc)
	rc.Close()
}
