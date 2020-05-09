package data

import (
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

// getVersion returns the SemVer version string from the tag in the format 0.0.0-rc0
func getVersion(s string) string {
	return versionReg.FindString(s)
}

func compareNumericString(a, b string) int {
	ai, err := strconv.Atoi(a)
	if err != nil {
		return -1
	}

	bi, err := strconv.Atoi(b)
	if err != nil {
		return 1
	}

	if ai < bi {
		return -1
	} else if ai > bi {
		return 1
	}

	return 0
}

// GetRC returns the rc portion of the version string
func GetRC(v string) string {
	sections := strings.Split(v, "-")

	if len(sections) <= 1 {
		return ""
	}

	return sections[1]
}

// VersionCmp compares two version strings
func VersionCmp(a, b string) int {
	l := 0

	a = versionReg.FindString(a)
	b = versionReg.FindString(b)

	aSections := strings.Split(a, "-")
	bSections := strings.Split(b, "-")
	aParts := strings.Split(aSections[0], ".")
	bParts := strings.Split(bSections[0], ".")
	aLen := len(aParts)
	bLen := len(bParts)

	if aLen > bLen {
		l = len(aParts)
	} else {
		l = len(bParts)
	}

	for i := 0; i < l; i++ {
		if i >= aLen {
			return -1
		} else if i >= bLen {
			return 1
		}

		v := compareNumericString(aParts[i], bParts[i])

		if v != 0 {
			return v
		}
	}

	aRC := len(aSections) != 1
	bRC := len(bSections) != 1

	if aRC && bRC {
		aRCStr := rcRegex.FindStringSubmatch(aSections[1])[0]
		bRCStr := rcRegex.FindStringSubmatch(bSections[1])[0]

		return compareNumericString(aRCStr, bRCStr)
	} else if aRC && !bRC {
		return -1
	} else if !aRC && bRC {
		return 1
	}

	return 0
}

func drainAndClose(rc io.ReadCloser) {
	io.Copy(ioutil.Discard, rc)
	rc.Close()
}
