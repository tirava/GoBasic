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
	Date    string // todo change to time.Time
	Summary string
	Body    template.HTML
}

// Posts storage todo - need xSQL storage instead map
var Posts map[string]Post

func init() {
	Posts = map[string]Post{
		"Мой первый пост!": { // todo may be convert to short name
			"Мой первый пост!",
			"18-е Сентября 2019 года",
			"Это короткое вступление.",
			template.HTML(blackfriday.Run([]byte(`
Здесь основной текст.
# Markdown!
*Это* **круто**!
`))),
		},
		"Это уже второй пост!": {
			"Это уже второй пост!",
			"19-е Сентября 2019 года",
			"Блог потихоньку растет.",
			template.HTML(blackfriday.Run([]byte(`
Разобрался в шаблонах и маркдаунах, как их совместить.

Теперь понять, как переходить на отдельные посты.
# Anybody!
*Hz* **cool**!
`))),
		},
	}
}

func mainPage(w http.ResponseWriter, _ *http.Request) {
	t := template.Must(template.ParseFiles(path.Join(templatePath, indexTemplate)))
	if err := t.Execute(w, Posts); err != nil {
		// todo buffer & return w code
		log.Println("error executing template in mainPage:", err)
	}
}

func postPage(w http.ResponseWriter, r *http.Request) {
	if _, ok := Posts[r.URL.Path[1:]]; !ok {
		log.Printf("post not found: %s", r.URL.Path[1:])
		w.WriteHeader(http.StatusNotFound)
		return
	}
	t := template.Must(template.ParseFiles(path.Join(templatePath, postTemplate)))
	if err := t.Execute(w, Posts[r.URL.Path[1:]]); err != nil {
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

	// prepare server
	mux := chi.NewRouter()
	mux.Get("/*", postPage)
	mux.Get("/", mainPage)

	fmt.Println("Starting server at:", servAddr)
	log.Fatalln(http.ListenAndServe(servAddr, mux))
}
