package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"snippetbox/cmd/web/config"
	"snippetbox/internal/models"
	"strconv"
)

func Home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			app.NotFound(w)
			return
		}
		files := []string{
			"./ui/html/base.tmpl.html",
			"./ui/html/pages/home.tmpl.html",
			"./ui/html/partials/nav.tmpl.html",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		err = ts.ExecuteTemplate(w, "base", nil)
		if err != nil {
			app.ServerError(w, err)
			return
		}
	}
}
func SnippetView(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			app.NotFound(w)
			return
		}

		snippet, err := app.Snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.NotFound(w)
			} else {
				app.ServerError(w, err)
			}
			return
		}

		fmt.Fprintf(w, "%+v", snippet)
	}
}
func SnippetCreate(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			app.ClientError(w, http.StatusMethodNotAllowed)
			return
		}

		title := "test"
		content := "test"
		expires := 7

		id, err := app.Snippets.Insert(title, content, expires)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating snippet: %v", err)
			app.ServerError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	}
}
