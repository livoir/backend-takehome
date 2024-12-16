package common

import "github.com/microcosm-cc/bluemonday"

var sanitizer *bluemonday.Policy

func Sanitize(input string) string {
	if sanitizer == nil {
		sanitizer = bluemonday.UGCPolicy()
	}
	return sanitizer.Sanitize(input)
}
