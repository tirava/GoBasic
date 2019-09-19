/*
 * HomeWork-3: Simple blog
 * Created on 19.09.2019 19:35
 * Copyright (c) 2019 - Eugene Klimov
 */
package main

import (
	"bytes"
	"github.com/go-chi/chi"
	"html/template"
	"log"
	"net/http"
	"path"
)

func mainPage(w http.ResponseWriter, _ *http.Request) {
	tIndex = template.Must(template.ParseFiles(path.Join(templatePath, indexTemplate))) // todo del
	var b bytes.Buffer                                                                  // no need to show bad content
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
	tPost = template.Must(template.ParseFiles(path.Join(templatePath, postTemplate))) // todo del
	postNum := chi.URLParam(r, "id")
	if _, ok := Posts[postNum]; !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	var b bytes.Buffer // no need to show bad content
	if err := tPost.ExecuteTemplate(&b, "post", Posts[postNum]); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if _, err := b.WriteTo(w); err != nil {
		log.Println("can't write to ResponseWriter in postPage")
	}
}
