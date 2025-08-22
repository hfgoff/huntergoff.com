package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"text/template"
	"time"

	"github.com/hfgoff/huntergoff/logging"
	"github.com/movieofthenight/go-streaming-availability/v4"
)

type config struct {
	port        int
	version     string
	commit      string
	rapidApikey string
}

type Result struct {
	Title            string
	Overview         string
	Poster           string
	ReleaseYear      int32
	StreamingOptions []streaming.StreamingOption
}

type PageData struct {
	Query   string
	Shows   []Result
	Version string
	Commit  string
}

var (
	streamingTmpl = template.Must(template.ParseFiles("templates/streaming.html"))
)

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "port")
	flag.StringVar(&cfg.version, "version", "(devel)", "version number")
	flag.StringVar(&cfg.rapidApikey, "rapidapikey", os.Getenv("RAPID_API_KEY"), "RapidAPI Key")

	flag.Parse()

	if info, ok := debug.ReadBuildInfo(); ok {
		cfg.version = "(devel)"
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" {
				cfg.version = s.Value
			}
		}
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/streaming", http.StatusFound)
	})

	mux.HandleFunc("/streaming", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")

		client := streaming.NewAPIClientFromRapidAPIKey(cfg.rapidApikey, nil)

		shows, _, err := client.ShowsAPI.
			SearchShowsByTitle(r.Context()).
			Title(query).
			Country("us").
			Execute()
		if err != nil {
			slog.Error("error searching shows", logging.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var page []Result
		for _, show := range shows {
			if _, ok := show.GetTitleOk(); !ok {
				continue
			} else if show.GetTitle() == "" {
				continue
			}
			page = append(page, Result{
				Title:            show.GetTitle(),
				Overview:         show.GetOverview(),
				Poster:           show.GetImageSet().VerticalPoster.W720,
				ReleaseYear:      show.GetReleaseYear(),
				StreamingOptions: show.GetStreamingOptions()["us"],
			})
		}

		data := PageData{
			Query:   query,
			Version: cfg.version,
			Commit:  cfg.commit,
			Shows:   page,
		}

		if err := streamingTmpl.Execute(w, data); err != nil {
			slog.Error("error executing streaming template", logging.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	srv := &http.Server{
		// strconv is faster than fmt.
		// see: https://github.com/uber-go/guide/blob/master/style.md#prefer-strconv-over-fmt
		Addr:         ":" + strconv.Itoa(cfg.port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("server error", logging.Error(err))
	}
}
