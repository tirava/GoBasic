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
	"strconv"
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
		log.Println("can't write to ResponseWriter in editPostPageForm")
	}
}

func (h *Handler) editPostPage(w http.ResponseWriter, r *http.Request) {
	postNum := chi.URLParam(r, "id")
	id, err := strconv.Atoi(postNum)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	post := Post{
		Id:      id,
		Title:   r.FormValue("title"),
		Date:    r.FormValue("date"),
		Summary: r.FormValue("summary"),
		Body:    template.HTML(r.FormValue("body")),
	}
	h.posts[postNum] = post
	http.Redirect(w, r, postsURL+"/"+postNum, http.StatusFound)
}

func (h *Handler) createPostPageForm(w http.ResponseWriter, r *http.Request) {
	h.tGlob = template.Must(template.ParseGlob(path.Join(templatePath, templateExt))) // todo del
	var b bytes.Buffer                                                                // no need to show bad content
	post := Post{}
	if err := h.tGlob.ExecuteTemplate(&b, "create", post); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if _, err := b.WriteTo(w); err != nil {
		log.Println("can't write to ResponseWriter in createPostPageForm")
	}
}

func (h *Handler) createPostPage(w http.ResponseWriter, r *http.Request) {
	id := len(h.posts) + 1
	h.posts[strconv.Itoa(id)] = Post{
		Id:      id,
		Title:   r.FormValue("title"),
		Date:    r.FormValue("date"),
		Summary: r.FormValue("summary"),
		Body:    template.HTML(r.FormValue("body")),
	}
	http.Redirect(w, r, postsURL+"/", http.StatusFound)
}

func (h *Handler) deletePostPage(w http.ResponseWriter, r *http.Request) {
	delete(h.posts, chi.URLParam(r, "id"))
	http.Redirect(w, r, postsURL+"/", http.StatusFound)
}

func (h *Handler) initPosts() {
	h.posts = dbPosts{
		"1": {
			Id:      1,
			Title:   "Мой первый пост!",
			Date:    "18-е Сентября 2019 года",
			Summary: "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Odio praesentium, quos. Aspernatur assumenda cupiditate deserunt ducimus, eveniet, expedita inventore laboriosam magni modi non odio, officia qui sequi similique unde voluptatem.",
			Body: `
Здесь основной текст.
# Markdown!
*Это* **круто**!
`,
		},
		"2": {
			Id:    2,
			Title: "Это уже второй пост!",
			Date:  "19-е Сентября 2019 года",
			Summary: `
Lorem ipsum dolor sit amet, consectetur adipisicing elit. Odio praesentium, quos. Aspernatur assumenda cupiditate deserunt ducimus, eveniet, expedita inventore laboriosam magni modi non odio, officia qui sequi similique unde voluptatem.
Lorem ipsum dolor sit amet, consectetur adipisicing elit. Odio praesentium, quos. Aspernatur assumenda cupiditate deserunt ducimus, eveniet, expedita inventore laboriosam magni modi non odio, officia qui sequi similique unde voluptatem.
`,
			Body: `
Разобрался в шаблонах и маркдаунах, как их совместить.

Теперь понять, как переходить на отдельные посты.
# Anybody!
*Hz* **cool**!
`,
		},
	}
}
