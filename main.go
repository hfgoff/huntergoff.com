package main

import (
	"net/http"
)

func main() {
	// This is a placeholder main function.
	// You can add your code here to run the application.
	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/", fs)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}
}
