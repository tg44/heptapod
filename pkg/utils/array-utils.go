package utils

import "os"

func ContainsFIA(s []os.FileInfo, e string) bool {
	for _, a := range s {
		if a.Name() == e {
			return true
		}
	}
	return false
}

func ContainsSA(a []string, s string) bool {
	for _, e := range a {
		if e == s {
			return true
		}
	}
	return false
}
