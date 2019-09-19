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
	tGlob = template.Must(template.ParseGlob(path.Join(templatePath, templateExt))) // todo del
	var b bytes.Buffer                                                              // no need to show bad content
	if err := tGlob.ExecuteTemplate(&b, "index", Posts); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if _, err := b.WriteTo(w); err != nil {
		log.Println("can't write to ResponseWriter in indexPage")
	}
}

func postPage(w http.ResponseWriter, r *http.Request) {
	tGlob = template.Must(template.ParseGlob(path.Join(templatePath, templateExt))) // todo del
	postNum := chi.URLParam(r, "id")
	if _, ok := Posts[postNum]; !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	var b bytes.Buffer // no need to show bad content
	if err := tGlob.ExecuteTemplate(&b, "post", Posts[postNum]); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if _, err := b.WriteTo(w); err != nil {
		log.Println("can't write to ResponseWriter in postPage")
	}
}
