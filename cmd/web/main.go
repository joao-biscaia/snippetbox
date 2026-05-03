package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"snippetbox/cmd/web/config"
	"snippetbox/internal/models"

	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	staticDir := flag.String("static-dir", "./ui/static", "Path to static assets")
	dsn := flag.String("dsn", "web:password@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()
	db, err := openDb(*dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	templateCache, err := config.NewTemplatecache()
	if err != nil {
		log.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	app := &config.Application{
		InfoLog:       log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog:      log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		Snippets:      &models.SnippetModel{DB: db},
		TemplateCache: templateCache,
		FormDecoder:   formDecoder,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.ErrorLog,
		Handler:  app.Routes(*staticDir),
	}

	app.InfoLog.Printf("Starting server on %v", srv.Addr)

	err = srv.ListenAndServe()
	app.ErrorLog.Fatal(err)
}

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
