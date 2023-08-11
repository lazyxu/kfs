package core

import "strings"

func FormatPath(p string) []string {
	splitPath := strings.Split(p, "/")
	retPath := make([]string, 0, len(splitPath))
	for i := 0; i < len(splitPath); i++ {
		if splitPath[i] != "" {
			retPath = append(retPath, splitPath[i])
		}
	}
	return retPath
}
