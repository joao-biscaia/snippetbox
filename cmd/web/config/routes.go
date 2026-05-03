package config

import (
	"net/http"
)

// The Routes() method returns a servemux containing our application Routes.
func (app *Application) Routes(staticDir string) http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", Home(app))
	mux.HandleFunc("/snippet/view", SnippetView(app))
	mux.HandleFunc("/snippet/create", SnippetCreate(app))

	return app.logRequest(SecureHeaders(mux))
}
