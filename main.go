package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/edgedb/edgedb-go"
)

type Application struct {
	DBClient *edgedb.Client
	CTX      context.Context
}

func main() {
	fmt.Println("START ")

	app := Application{}

	ctx := context.Background()
	app.CTX = ctx
	opts := edgedb.Options{
		Database:    "main",
		User:        "edgedb",
		Concurrency: 4,
	}
	client, err := edgedb.CreateClient(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	app.DBClient = client

	mux := http.NewServeMux()
	mux.HandleFunc("GET /movies", app.getMovies)

	http.ListenAndServe(":6565", mux)
}

type Movie struct {
	ID    edgedb.UUID `edgedb:"id"`
	Title string      `edgedb:"title"`
}

func (app *Application) getMovies(w http.ResponseWriter, r *http.Request) {

	movies := []Movie{}

	query := "SELECT Movie{id,title}"
	err := app.DBClient.Query(app.CTX, query, &movies)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%+v", movies)

	err = app.writeJSON(w, http.StatusOK, envelope{"movies": movies}, nil)
	if err != nil {
		log.Fatal(err)
	}

}

type envelope map[string]any

func (app *Application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}
