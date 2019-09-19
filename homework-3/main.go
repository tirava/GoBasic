/*
 * HomeWork-3: Simple blog
 * Created on 18.09.2019 23:11
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

const (
	servAddr      = "localhost:8080"
	indexTemplate = "index.gohtml"
	postTemplate  = "post.gohtml"
	templatePath  = "templates"
	postsURL      = "/posts"
)

func main() {

	// prepare routes & middleware
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Get("/", mainPage)
	//mux.Get(postsURL, mainPage)
	//mux.Get(postsURL+"/{id}", postPage)

	mux.Route(postsURL, func(r chi.Router) {
		r.Get("/{id}", postPage)
		r.Get("/", mainPage)
		//r.Get("/new", newPostPage)
		//r.Post("/new", newPostPage)
	})

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "static")
	mux.Handle("/static", http.FileServer(http.Dir(filesDir)))

	srv := &http.Server{Addr: servAddr, Handler: mux}

	// graceful shutdown
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM) // os.Kill cannot be trapped anyway!

	go func() {
		log.Println("Signal received:", <-shutdown)
		if err := srv.Shutdown(nil); err != nil {
			log.Println("Error while shutdown server:", err)
		}
	}()

	fmt.Println("Starting server at:", servAddr)
	log.Printf("Shutdown server at: %s\n%v", servAddr, srv.ListenAndServe())
}
