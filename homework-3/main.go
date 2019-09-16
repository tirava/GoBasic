/*
 * HomeWork-3: Simple blog
 * Created on 18.09.2019 23:11
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	servAddr      = "localhost:8080"
	indexTemplate = "templates/index.gohtml"
	postsPath     = "posts/"
)

type Post struct {
	Title   string
	Date    string
	Summary string
	Body    string
	File    string
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	posts := getPosts()
	t := template.Must(template.ParseFiles(indexTemplate))
	t.Execute(w, posts)
}

func getPosts() []Post {
	posts := &[]Post{}
	files, _ := filepath.Glob(postsPath + "*")
	for _, f := range files {
		file := strings.ReplaceAll(f, postsPath, "") // remove path
		file = strings.ReplaceAll(file, ".md", "")   // remove ext
		rf, _ := ioutil.ReadFile(f)
		lines := strings.Split(string(rf), "\n")
		title := string(lines[0])
		date := string(lines[1])
		summary := string(lines[2])
		body := strings.Join(lines[3:len(lines)], "\n")
		body = string(blackfriday.MarkdownCommon([]byte(body)))
		*posts = append(*posts, Post{title, date, summary, body, file})
	}
	return *posts
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
	http.HandleFunc("/", mainPage)

	fmt.Println("Starting server at:", servAddr)
	log.Fatalln(http.ListenAndServe(servAddr, nil))
}
