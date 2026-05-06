package config

import (
	"html/template"
	"log"
	"snippetbox/internal/models"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
)

type Application struct {
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	Snippets       *models.SnippetModel
	Users          *models.UserModel
	TemplateCache  map[string]*template.Template
	FormDecoder    *form.Decoder
	SessionManager *scs.SessionManager
}
