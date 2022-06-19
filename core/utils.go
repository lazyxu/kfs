package core

import "strings"

func FormatPath(p string) []string {
	splitPath := strings.Split(p, "/")
	retPath := make([]string, 0, len(splitPath))
	index := 0
	for i := 0; i < len(splitPath); i++ {
		if splitPath[i] != "" {
			retPath[index] = splitPath[i]
			index++
		}
	}
	return retPath
}
