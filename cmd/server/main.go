package main

import (
	"errors"
	"log"
	"net/http"

	"website-monitor/internal/app"
	"website-monitor/internal/storage"
)

func main() {
	store := storage.NewStorage()

	store.Add(storage.Site{URL: "https://example.com"})

	router := app.NewRouter(store)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

}
