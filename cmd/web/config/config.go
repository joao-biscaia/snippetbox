package config

import (
	"log"
	"net/http"
	"snippetbox/internal/models"
)

type Application struct {
	ErrorLog    *log.Logger
	InfoLog     *log.Logger
	Snippets    *models.SnippetModel
	ServerError func(w http.ResponseWriter, err error)
	ClientError func(w http.ResponseWriter, status int)
	NotFound    func(w http.ResponseWriter)
}
