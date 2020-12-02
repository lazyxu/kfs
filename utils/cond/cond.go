package cond

func String(b bool, x string, y string) string {
	if b {
		return x
	}
	return y
}
