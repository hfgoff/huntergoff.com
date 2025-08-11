package main

import (
	"cmp"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/hfgoff/huntergoff/logging"
)

type Config struct {
	Port int
}

type PageData struct {
	Title   string
	Version string
	Commit  string
}

type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

var (
	version = "unknown"
	commit  = "?"
	date    = ""
)

func main() {
	cfg := Config{
		Port: 8080,
	}
	info := createBuildInfo()

	logger := logging.New("huntergoff-com").With(slog.String("version", info.Version))

	flag.IntVar(&cfg.Port, "port", cfg.Port, "Port to run the server on")
	flag.Parse()

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		logger.Error("Failed to parse templates", "error", err)
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Title:   "Hunter Goff",
			Version: info.Version,
			Commit:  info.Commit,
		}

		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("Failed to execute template", "error", err)
			return
		}
	})

	logger.Info("Starting server", "port", cfg.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
	if err != nil {
		panic(err)
	}
}

func createBuildInfo() BuildInfo {
	info := BuildInfo{
		Commit:  commit,
		Version: version,
		Date:    date,
	}

	buildInfo, available := debug.ReadBuildInfo()
	if !available {
		return info
	}

	if date != "" {
		return info
	}

	info.Version = buildInfo.Main.Version

	var revision string
	var modified string
	for _, setting := range buildInfo.Settings {
		// The `vcs.xxx` information is only available with `go build`.
		// This information is not available with `go install` or `go run`.
		switch setting.Key {
		case "vcs.time":
			info.Date = setting.Value
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			modified = setting.Value
		}
	}

	info.Date = cmp.Or(info.Date, "(unknown)")

	info.Commit = fmt.Sprintf("(%s, modified: %s)",
		cmp.Or(revision, "unknown"), cmp.Or(modified, "?"))

	return info
}
