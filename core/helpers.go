package core

func SliceContains(name string, list []string) bool {
	for _, b := range list {
		if b == name {
			return true
		}
	}
	return false
}
