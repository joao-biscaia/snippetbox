package config

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"snippetbox/internal/models"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.Snippets.Latest()
	if err != nil {
		app.ServerError(w, err)
		return
	}

	data := app.NewTemplateData(r)
	data.Snippets = snippets

	app.Render(w, http.StatusOK, "home.tmpl.html", data)
}
func (app *Application) SnippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
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

	data := app.NewTemplateData(r)
	data.Snippet = snippet

	app.Render(w, http.StatusOK, "view.tmpl.html", data)
}
func (app *Application) SnippetCreatePost(w http.ResponseWriter, r *http.Request) {
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

func (app *Application) SnippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet..."))
}
