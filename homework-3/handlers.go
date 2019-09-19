/*
 * HomeWork-3: Simple blog
 * Created on 19.09.2019 19:35
 * Copyright (c) 2019 - Eugene Klimov
 */
package main

import (
	"bytes"
	"github.com/go-chi/chi"
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"log"
	"net/http"
	"path"
)

func (h *Handler) mainPage(w http.ResponseWriter, _ *http.Request) {
	h.tGlob = template.Must(template.ParseGlob(path.Join(templatePath, templateExt))) // todo del
	var b bytes.Buffer                                                                // no need to show bad content
	if err := h.tGlob.ExecuteTemplate(&b, "index", h.posts); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if _, err := b.WriteTo(w); err != nil {
		log.Println("can't write to ResponseWriter in indexPage")
	}
}

func (h *Handler) postPage(w http.ResponseWriter, r *http.Request) {
	h.tGlob = template.Must(template.ParseGlob(path.Join(templatePath, templateExt))) // todo del
	postNum := chi.URLParam(r, "id")
	if _, ok := h.posts[postNum]; !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	post := h.posts[postNum] // create copy for markdown convert
	post.Body = template.HTML(blackfriday.Run([]byte(post.Body)))
	var b bytes.Buffer // no need to show bad content
	if err := h.tGlob.ExecuteTemplate(&b, "post", post); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if _, err := b.WriteTo(w); err != nil {
		log.Println("can't write to ResponseWriter in postPage")
	}
}

func (h *Handler) editPostPageForm(w http.ResponseWriter, r *http.Request) {
	h.tGlob = template.Must(template.ParseGlob(path.Join(templatePath, templateExt))) // todo del
	postNum := chi.URLParam(r, "id")
	var b bytes.Buffer // no need to show bad content
	if err := h.tGlob.ExecuteTemplate(&b, "edit", h.posts[postNum]); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if _, err := b.WriteTo(w); err != nil {
		log.Println("can't write to ResponseWriter in newPostPageForm")
	}
}

func (h *Handler) editPostPage(w http.ResponseWriter, r *http.Request) {
	//h.posts
}

func (h *Handler) initPosts() {
	h.posts = dbPosts{
		"1": {
			1,
			"Мой первый пост!",
			"18-е Сентября 2019 года",
			"Lorem ipsum dolor sit amet, consectetur adipisicing elit. Odio praesentium, quos. Aspernatur assumenda cupiditate deserunt ducimus, eveniet, expedita inventore laboriosam magni modi non odio, officia qui sequi similique unde voluptatem.",
			`
Здесь основной текст.
# Markdown!
*Это* **круто**!
`,
		},
		"2": {
			2,
			"Это уже второй пост!",
			"19-е Сентября 2019 года",
			`
Lorem ipsum dolor sit amet, consectetur adipisicing elit. Odio praesentium, quos. Aspernatur assumenda cupiditate deserunt ducimus, eveniet, expedita inventore laboriosam magni modi non odio, officia qui sequi similique unde voluptatem.
Lorem ipsum dolor sit amet, consectetur adipisicing elit. Odio praesentium, quos. Aspernatur assumenda cupiditate deserunt ducimus, eveniet, expedita inventore laboriosam magni modi non odio, officia qui sequi similique unde voluptatem.
`,
			`
Разобрался в шаблонах и маркдаунах, как их совместить.

Теперь понять, как переходить на отдельные посты.
# Anybody!
*Hz* **cool**!
`,
		},
	}
}
