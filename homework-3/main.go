/*
 * HomeWork-3: Simple blog
 * Created on 18.09.2019 23:11
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	servAddr      = "localhost:8080"
	indexTemplate = "index.gohtml"
	postTemplate  = "post.gohtml"
	templatePath  = "templates" // + string(os.PathSeparator)
	postsPath     = "posts"     // + string(os.PathSeparator)
	postsExt      = ".md"
)

type Post struct {
	Title   string
	Date    string
	Summary string
	Body    template.HTML
	File    string
}

func mainPage(w http.ResponseWriter, _ *http.Request) {
	posts, err := getPosts()
	if err != nil {
		log.Println("error getting posts:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t := template.Must(template.ParseFiles(path.Join(templatePath, indexTemplate)))
	err = t.Execute(w, posts)
	if err != nil {
		// todo buffer & return w code
	}
	return
}

func postPage(w http.ResponseWriter, r *http.Request) {
	f := path.Join(postsPath, r.URL.Path[1:]) + postsExt
	rf, err := ioutil.ReadFile(f)
	if err != nil {
		log.Println("error reading file:", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	post := lines2Posts(rf, r.URL.Path[1:])
	t := template.Must(template.ParseFiles(path.Join(templatePath, postTemplate)))
	err = t.Execute(w, post)
	if err != nil {
		// todo buffer & return w code
		log.Println(err)
	}
}

func getPosts() ([]Post, error) {
	var posts []Post
	files, err := filepath.Glob(path.Join(postsPath, "*") + postsExt)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		file := strings.ReplaceAll(f, postsPath+string(os.PathSeparator), "") // remove path
		file = strings.ReplaceAll(file, postsExt, "")                         // remove ext
		rf, err := ioutil.ReadFile(f)
		if err != nil {
			continue // skip bad file
		}
		posts = append(posts, lines2Posts(rf, file))
	}
	return posts, nil
}

func lines2Posts(readFile []byte, file string) Post {
	lines := strings.Split(string(readFile), "\n")
	title, date, summary := lines[0], lines[1], lines[2]
	body := template.HTML(strings.Join(lines[3:], "\n"))
	post := Post{
		title, date, summary,
		template.HTML(blackfriday.Run([]byte(body))),
		file,
	}
	return post
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
