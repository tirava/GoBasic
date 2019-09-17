/*
 * HomeWork-3: Simple blog
 * Created on 18.09.2019 23:11
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"fmt"
	"github.com/go-chi/chi"
	//"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
)

const (
	servAddr      = "localhost:8080"
	indexTemplate = "index.gohtml"
	postTemplate  = "post.gohtml"
	templatePath  = "templates"
	//postsPath     = "posts"
	//postsExt      = ".md"
)

// Post is the base post type
type Post struct {
	Title   string
	Date    string
	Summary string
	Body    template.HTML
}

// Posts storage todo - need SQL storage instead map
var Posts map[string]Post

func init() {
	Posts = map[string]Post{
		"111": {
			"111",
			"222",
			"333",
			"444",
		},
		"555": {
			"555",
			"666",
			"777",
			"888",
		},
	}
}

func mainPage(w http.ResponseWriter, _ *http.Request) {
	t := template.Must(template.ParseFiles(path.Join(templatePath, indexTemplate)))
	err := t.Execute(w, Posts)
	if err != nil {
		// todo buffer & return w code
		log.Println("error executing template in mainPage:", err)
	}
	return
}

func postPage(w http.ResponseWriter, r *http.Request) {
	if _, ok := Posts[r.URL.Path[1:]]; !ok {
		log.Printf("post not found: %s", r.URL.Path[1:])
		w.WriteHeader(http.StatusNotFound)
		return
	}
	t := template.Must(template.ParseFiles(path.Join(templatePath, postTemplate)))
	err := t.Execute(w, Posts[r.URL.Path[1:]])
	if err != nil {
		// todo buffer & return w code
		log.Println("error executing template in postPage:", err)
	}
}

func main() {

	// safe shutdown
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM) // os.Kill wrong for linters

	go func() {
		<-shutdown
		// any work here
		fmt.Printf("\nShutdown server at: %s\n", servAddr)
		os.Exit(0)
	}()

	// prepare server, no need smart router for simple scenario
	mux := chi.NewRouter()
	mux.Get("/*", postPage)
	mux.Get("/", mainPage)

	fmt.Println("Starting server at:", servAddr)
	log.Fatalln(http.ListenAndServe(servAddr, mux))
}
