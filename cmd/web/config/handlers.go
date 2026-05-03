package config

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/internal/models"
	"snippetbox/internal/validator"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type SnippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
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
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank.")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365.")

	if !form.Valid() {
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
