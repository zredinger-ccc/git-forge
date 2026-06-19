package version

import "runtime/debug"

var (
	Version = "dev"
	Commit  = ""
	Date    = ""
)

func Info() map[string]string {
	commit := Commit
	date := Date
	if commit == "" || date == "" {
		if bi, ok := debug.ReadBuildInfo(); ok {
			for _, s := range bi.Settings {
				switch s.Key {
				case "vcs.revision":
					if commit == "" {
						commit = s.Value
					}
				case "vcs.time":
					if date == "" {
						date = s.Value
					}
				}
			}
		}
	}
	return map[string]string{
		"version": Version,
		"commit":  commit,
		"date":    date,
	}
}
