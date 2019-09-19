/*
 * HomeWork-3: Simple blog
 * Created on 18.09.2019 23:11
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"bytes"
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
	"syscall"
	"time"
)

const (
	servAddr      = "localhost:8080"
	indexTemplate = "index.gohtml"
	postTemplate  = "post.gohtml"
	templatePath  = "templates"
	postsURL      = "/posts"
)

// Post is the base post type
type Post struct {
	Title   string
	Date    string // todo change to time.Time
	Summary string
	Body    template.HTML
}

func mainPage(w http.ResponseWriter, _ *http.Request) {
	tIndex = template.Must(template.ParseFiles(path.Join(templatePath, indexTemplate)))
	var b bytes.Buffer // no need to show bad content
	//if err := tIndex.Execute(&b, Posts); err != nil {
	if err := tIndex.ExecuteTemplate(&b, "index", Posts); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if _, err := b.WriteTo(w); err != nil {
		log.Println("can't write to ResponseWriter in indexPage")
	}
}

func postPage(w http.ResponseWriter, r *http.Request) {
	tPost = template.Must(template.ParseFiles(path.Join(templatePath, postTemplate)))
	postNum := chi.URLParam(r, "id")
	if _, ok := Posts[postNum]; !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	var b bytes.Buffer // no need to show bad content
	if err := tPost.Execute(&b, Posts[postNum]); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if _, err := b.WriteTo(w); err != nil {
		log.Println("can't write to ResponseWriter in postPage")
	}
}

func main() {

	// prepare server
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Get("/", mainPage)
	mux.Get(postsURL, mainPage)
	mux.Get(postsURL+"/{id}", postPage)
	//mux.Get("/posts/new", newPostPage)
	//mux.Post("/posts/new", newPostPage)

	srv := &http.Server{
		Addr:    servAddr,
		Handler: mux,
	}

	// graceful shutdown
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM) // os.Kill cannot be trapped anyway!
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		log.Println("Signal received:", <-shutdown)
		// any work here
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalln("Error while shutdown server:", err)
		}
	}()

	fmt.Println("Starting server at:", servAddr)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Shutdown server at: %s\n%v", servAddr, err)
	}
}
