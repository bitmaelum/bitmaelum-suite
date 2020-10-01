// +build !windows

package config

func getSearchPaths() []string {
	return []string{
		"./",
		"/etc/bitmaelum/",
	}
}
