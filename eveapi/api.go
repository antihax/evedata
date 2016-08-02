package eveapi

import "regexp"

func IsValidVCode(vc string) bool {
	if m, _ := regexp.MatchString("^[a-zA-Z0-9]+$", vc); !m {
		return false
	}

	if len(vc) != 64 {
		return false
	}

	return true
}
