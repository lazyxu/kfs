package localfs

func contains(array []string, s string) bool {
	for _, elm := range array {
		if elm == s {
			return true
		}
	}
	return false
}
