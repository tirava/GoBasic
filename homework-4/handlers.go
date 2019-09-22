/*
 * HomeWork-4: Simple blog - MySQL
 * Created on 22.09.2019 13:11
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"github.com/go-chi/chi"
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) mainPageForm(w http.ResponseWriter, _ *http.Request) {
	if err := h.tmplGlob.ExecuteTemplate(w, "index", h.posts); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) postsPageForm(w http.ResponseWriter, r *http.Request) {
	postNum := r.URL.Query().Get("id")
	if postNum == "" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
	if _, ok := h.posts[postNum]; !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	post := h.posts[postNum] // create copy for markdown convert
	post.Body = template.HTML(blackfriday.Run([]byte(post.Body)))
	if err := h.tmplGlob.ExecuteTemplate(w, "post", post); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) editPostPageForm(w http.ResponseWriter, r *http.Request) {
	postNum := r.URL.Query().Get("id")
	if err := h.tmplGlob.ExecuteTemplate(w, "edit", h.posts[postNum]); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
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
		ID:      id,
		Title:   r.FormValue(TITLE),
		Date:    r.FormValue(DATE),
		Summary: r.FormValue(SUMMARY),
		Body:    template.HTML(r.FormValue("body")),
	}
	h.posts[postNum] = post
	http.Redirect(w, r, POSTSURL+"/?id="+postNum, http.StatusFound)
}

func (h *Handler) createPostPageForm(w http.ResponseWriter, r *http.Request) {
	post := Post{}
	if err := h.tmplGlob.ExecuteTemplate(w, "create", post); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) createPostPage(w http.ResponseWriter, r *http.Request) {
	id := h.nextGlobID()
	ids := strconv.Itoa(id)
	h.posts[ids] = Post{
		ID:      id,
		Title:   r.FormValue(TITLE),
		Date:    r.FormValue(DATE),
		Summary: r.FormValue(SUMMARY),
		Body:    template.HTML(r.FormValue("body")),
	}
	http.Redirect(w, r, POSTSURL+"/?id="+ids, http.StatusFound)
}

func (h *Handler) deletePostPage(w http.ResponseWriter, r *http.Request) {
	delete(h.posts, chi.URLParam(r, "id"))
	http.Redirect(w, r, POSTSURL+"/", http.StatusFound)
}

func (h *Handler) nextGlobID() int {
	h.mux.Lock()
	h.globID++
	h.mux.Unlock()
	return h.globID
}

func (h *Handler) initPosts() {
	h.posts = dbPosts{
		"1": {
			ID:      1,
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
			ID:    2,
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
		"3": {
			ID:      3,
			Title:   "Пример основных вариантов разметки Markdown",
			Date:    "20-е Сентября 2019 года",
			Summary: "Официальное руководство по синтаксису Markdown мне кажется слишком длинным и не слишком наглядным, поэтому я составил краткое руководство, которое поможет выучить или повторить синтаксис Маркдауна за час.",
			Body:    `todo`,
		},
	}
	h.globID = 3 // last id
}
