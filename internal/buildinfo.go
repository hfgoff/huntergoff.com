package internal

import "runtime/debug"

func Commit() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" {
				return s.Value // full git commit hash
			}
		}
	}
	return "(unknown)"
}
