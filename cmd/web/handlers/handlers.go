package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"snippetbox/cmd/web/config"
	"snippetbox/cmd/web/helpers"
	"strconv"
)

func Home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			helpers.NotFound(app)(w)
			return
		}
		files := []string{
			"./ui/html/base.tmpl.html",
			"./ui/html/pages/home.tmpl.html",
			"./ui/html/partials/nav.tmpl.html",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			helpers.ServerError(app)(w, err)
			return
		}

		err = ts.ExecuteTemplate(w, "base", nil)
		if err != nil {
			helpers.ServerError(app)(w, err)
			return
		}
	}
}
func SnippetView(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			helpers.NotFound(app)(w)
			return
		}
		fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	}
}
func SnippetCreate(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			helpers.ClientError(app)(w, http.StatusMethodNotAllowed)
			return
		}
		w.Write([]byte("Create a new snippet..."))
	}
}
