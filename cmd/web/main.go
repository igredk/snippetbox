package main

import (
	"context"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/igredk/snippetbox/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	pgDsn := flag.String("pgDsn", "postgres://biba:boba@localhost:5432/snippetbox", "Postgres connection URL")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dbPool, err := openDB(*pgDsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer dbPool.Close()

	templateCache, err := newTemplateCache()
    if err != nil {
        errorLog.Fatal(err)
    }

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: dbPool},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}
