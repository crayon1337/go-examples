package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	router := mux.NewRouter()

	bookRouter := router.PathPrefix("/books").Subrouter()

	bookRouter.HandleFunc("/index", AllBooks)
	bookRouter.HandleFunc("/create", CreateBook).Methods("POST")

	server := http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("%s\n", err)
		}
	}()

	log.Printf("Server started on %s\n", server.Addr)

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)

	defer cancel()

	server.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}

type Book struct {
	Title   string
	Content string
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Book creation %+v\n", book)
	fmt.Fprintf(w, "Book title: %s\n", book.Title)
}

func AllBooks(w http.ResponseWriter, r *http.Request) {

}
