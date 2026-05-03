package config

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/internal/models"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
)

type SnippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

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
	// Parses all form data into `PostForm` map inside http.Request param
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	form := &SnippetCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be longer than 100 characters"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be empty"
	}

	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "This field must equal do 1, 7 or 365"
	}

	if len(form.FieldErrors) > 0 {
		data := app.NewTemplateData(r)
		data.Form = form
		app.Render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.Snippets.Insert(form.Title, form.Content, expires)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *Application) SnippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)

	data.Form = &SnippetCreateForm{
		Expires: 365,
	}

	app.Render(w, http.StatusOK, "create.tmpl.html", data)

}
