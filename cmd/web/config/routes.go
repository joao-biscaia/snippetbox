package config

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// The Routes() method returns a servemux containing our application Routes.
func (app *Application) Routes(staticDir string) http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NotFound(w)
	})

	fileServer := http.FileServer(http.Dir(staticDir))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.Home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.SnippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.SnippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.SnippetCreatePost)

	standard := alice.New(app.recoverPanic, app.logRequest, SecureHeaders)

	return standard.Then(router)
}
