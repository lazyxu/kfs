package core

import "strings"

func FormatPath(p string) []string {
	splitPath := strings.Split(p, "/")
	if splitPath[0] == "" {
		splitPath = splitPath[1:]
	}
	return splitPath
}
