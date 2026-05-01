package main

import (
	"net/http"
	"snippetbox/cmd/web/config"
	"snippetbox/cmd/web/handlers"
)

// The routes() method returns a servemux containing our application routes.
func routes(app *config.Application, staticDir string) *http.ServeMux {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", handlers.Home(app))
	mux.HandleFunc("/snippet/view", handlers.SnippetView(app))
	mux.HandleFunc("/snippet/create", handlers.SnippetCreate(app))

	return mux
}
