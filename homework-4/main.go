/*
 * HomeWork-4: Simple blog - MySQL
 * Created on 22.09.2019 13:11
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

// Constants
const (
	SERVADDR     = ":8080"
	TEMPLATEEXT  = "*.gohtml"
	TEMPLATEPATH = "templates"
	POSTSURL     = "/posts"
	EDITURL      = "/edit"
	DELETEURL    = "/delete"
	CREATEURL    = "/create"
	APIURL       = "/api/v1"
	STATICPATH   = "/static"
	TITLE        = "title"
	DATE         = "date"
	SUMMARY      = "summary"
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
		tmplGlob: template.Must(template.ParseGlob(path.Join(TEMPLATEPATH, TEMPLATEEXT))),
	}
	handlers.initPosts()

	// prepare routes & middleware
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Get("/", handlers.mainPageForm)
	mux.Route(APIURL, func(r chi.Router) {
		r.Route(POSTSURL, func(r chi.Router) {
			r.Post(CREATEURL, handlers.createPostPage)
			r.Post(EDITURL+"/{id}", handlers.editPostPage)
			r.Post(DELETEURL+"/{id}", handlers.deletePostPage)
		})
	})
	mux.Route(POSTSURL, func(r chi.Router) {
		r.Get("/", handlers.postsPageForm)
		r.Get(CREATEURL, handlers.createPostPageForm)
		r.Get(EDITURL+"/", handlers.editPostPageForm)
	})

	// custom server is for custom parameters & graceful shutdown
	srv := &http.Server{Addr: SERVADDR, Handler: mux}

	// static files
	mux.Handle(STATICPATH+"/*", http.StripPrefix(STATICPATH, http.FileServer(http.Dir("."+STATICPATH))))

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

	fmt.Println("Starting server at:", SERVADDR)
	log.Printf("Shutdown server at: %s\n%v", SERVADDR, srv.ListenAndServe())
}
