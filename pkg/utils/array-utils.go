package utils

import "os"

func ContainsFIA(s []os.DirEntry, e string) bool {
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

func Unique(arr []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range arr {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// based on https://stackoverflow.com/a/48123201
func Filter(arr []string, f func(string) bool) []string {
	j := 0
	q := make([]string, len(arr))
	for _, n := range arr {
		if f(n) {
			q[j] = n
			j++
		}
	}
	q = q[:j]
	return q
}
