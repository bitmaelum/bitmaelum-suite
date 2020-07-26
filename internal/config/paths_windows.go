// +build windows

package config

func getSearchPaths() []string {
	return []string{
		"./",
		"%ProgramData%/BitMaelum/etc/bitmaelum",
	}
}
