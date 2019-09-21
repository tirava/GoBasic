/*
 * HomeWork-3: Simple blog
 * Created on 18.09.2019 23:11
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
)

const (
	servAddr     = ":8080"
	templateExt  = "*.gohtml"
	templatePath = "templates"
	postsURL     = "/posts"
	editURL      = "/edit"
	deleteURL    = "/delete"
	createURL    = "/create"
	apiURL       = "/api/v1"
	staticPath   = "/static"
)

// Post is the base post type
type Post struct {
	ID      int
	Title   string
	Date    string // todo change to time.Time
	Summary string
	Body    template.HTML
}

type dbPosts map[string]Post

// Handler is the global handlers struct
type Handler struct {
	posts    dbPosts
	tmplGlob *template.Template
	globID   int
	mux      sync.Mutex
}

func main() {

	// new handler struct
	handlers := &Handler{
		tmplGlob: template.Must(template.ParseGlob(path.Join(templatePath, templateExt))),
	}
	handlers.initPosts()

	// prepare routes & middleware
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Get("/", handlers.mainPageForm)
	mux.Route(apiURL, func(r chi.Router) {
		r.Route(postsURL, func(r chi.Router) {
			r.Post(createURL, handlers.createPostPage)
			r.Post(editURL+"/{id}", handlers.editPostPage)
			r.Post(deleteURL+"/{id}", handlers.deletePostPage)
		})
	})
	mux.Route(postsURL, func(r chi.Router) {
		r.Get("/", handlers.mainPageForm)
		r.Get("/{id}", handlers.postPageForm)
		r.Get(createURL, handlers.createPostPageForm)
		r.Get(editURL+"/{id}", handlers.editPostPageForm)
	})

	// custom server is for custom parameters & graceful shutdown
	srv := &http.Server{Addr: servAddr, Handler: mux}

	// static files
	mux.Handle(staticPath+"/*", http.StripPrefix(staticPath, http.FileServer(http.Dir("."+staticPath))))

	// graceful shutdown
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM) // os.Kill cannot be trapped anyway!
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		log.Println("Signal received:", <-shutdown)
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("Error while shutdown server:", err)
		}
	}()

	fmt.Println("Starting server at:", servAddr)
	log.Printf("Shutdown server at: %s\n%v", servAddr, srv.ListenAndServe())
}
