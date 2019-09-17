/*
 * HomeWork-3: Simple blog
 * Created on 18.09.2019 23:11
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"fmt"
	"gopkg.in/russross/blackfriday.v2"
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
	postTemplate  = "templates/post.gohtml"
	postsPath     = "posts/"
	postsExt      = ".md"
)

type Post struct {
	Title   string
	Date    string
	Summary string
	Body    string
	File    string
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[1:] == "" { // root handler
		posts, err := getPosts()
		if err != nil {
			// todo - return w code 500
		}
		t := template.Must(template.ParseFiles(indexTemplate))
		//
		err = t.Execute(w, posts)
		if err != nil {
			// todo buffer
		}
		return
	}
	f := postsPath + r.URL.Path[1:] + postsExt
	rf, err := ioutil.ReadFile(f)
	if err != nil {
		//todo return error for bad file? & return w code
		log.Println(err)
		return
	}
	lines := strings.Split(string(rf), "\n")
	title := string(lines[0])
	date := string(lines[1])
	summary := string(lines[2])
	body := strings.Join(lines[3:], "\n")
	post := Post{title, date, summary, body, r.URL.Path[1:]}

	t := template.Must(template.New("post.gohtml").Funcs(template.FuncMap{"markDown": markDowner}).ParseFiles(postTemplate))
	err = t.Execute(w, post)

	if err != nil {
		// todo buffer & return w code
		log.Println(err)
	}
}

func markDowner(args ...interface{}) template.HTML {
	s := blackfriday.Run([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}

func getPosts() ([]Post, error) {
	posts := &[]Post{}
	files, err := filepath.Glob(postsPath + "*" + postsExt)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		file := strings.ReplaceAll(f, postsPath, "")  // remove path
		file = strings.ReplaceAll(file, postsExt, "") // remove ext
		rf, err := ioutil.ReadFile(f)
		if err != nil {
			continue // skip bad file
		}
		lines := strings.Split(string(rf), "\n")
		title, date, summary := lines[0], lines[1], lines[2]
		body := strings.Join(lines[3:], "\n")
		//body = string(blackfriday.Run([]byte(body)))
		*posts = append(*posts, Post{title, date, summary, body, file})
	}
	return *posts, nil
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
