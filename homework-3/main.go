/*
 * HomeWork-3: Simple blog
 * Created on 18.09.2019 23:11
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"bytes"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
)

const (
	servAddr      = "localhost:8080"
	indexTemplate = "index.gohtml"
	postTemplate  = "post.gohtml"
	templatePath  = "templates"
	postsURL      = "posts/"
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

var tIndex, tPost *template.Template

func init() {
	tIndex = template.Must(template.ParseFiles(path.Join(templatePath, indexTemplate)))
	tPost = template.Must(template.ParseFiles(path.Join(templatePath, postTemplate)))

	Posts = map[string]Post{
		"1": {
			"Мой первый пост!",
			"18-е Сентября 2019 года",
			"Это короткое вступление.",
			template.HTML(blackfriday.Run([]byte(`
Здесь основной текст.
# Markdown!
*Это* **круто**!
`))),
		},
		"2": {
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
	var b bytes.Buffer // no need to show bad content
	if err := tIndex.Execute(&b, Posts); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if _, err := b.WriteTo(w); err != nil {
		log.Println("can't write to ResponseWriter in indexPage")
	}
}

func postPage(w http.ResponseWriter, r *http.Request) {
	postNum := strings.Replace(r.URL.Path[1:], postsURL, "", 1)
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
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Get("/*", postPage)
	mux.Get("/", mainPage)

	fmt.Println("Starting server at:", servAddr)
	log.Fatalln(http.ListenAndServe(servAddr, mux))
}
