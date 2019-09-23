/*
 * HomeWork-4: Simple blog - MySQL
 * Created on 22.09.2019 13:11
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// main page
func (h *Handler) mainPageForm(w http.ResponseWriter, _ *http.Request) {
	if err := h.tmplGlob.ExecuteTemplate(w, "index", h.posts); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// one post page
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
	h.execTemplate(w, post, "post")
}

// crete post page
func (h *Handler) createPostPageForm(w http.ResponseWriter, r *http.Request) {
	h.execTemplate(w, Post{}, "create")
}

// edit post page
func (h *Handler) editPostPageForm(w http.ResponseWriter, r *http.Request) {
	post := h.posts[r.URL.Query().Get("id")]
	h.execTemplate(w, post, "edit")
}

// exec template helper
func (h *Handler) execTemplate(w http.ResponseWriter, post Post, tmpl string) {
	if err := h.tmplGlob.ExecuteTemplate(w, tmpl, post); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// api create post
func (h *Handler) createPostPage(w http.ResponseWriter, r *http.Request) {
	post := &Post{}
	err := json.NewDecoder(r.Body).Decode(post)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err, "error while decoding post body for create")
		return
	}
	post.ID = h.nextGlobID()
	if err := h.posts.create(post); err != nil {
		h.sendError(w, http.StatusInternalServerError, err, "error while create post")
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// api edit post
func (h *Handler) editPostPage(w http.ResponseWriter, r *http.Request) {
	postNum := chi.URLParam(r, "id")
	post := &Post{}
	err := json.NewDecoder(r.Body).Decode(post)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err, "error while decoding post body for update")
		return
	}
	id, err := strconv.Atoi(postNum)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err, "error while preparing post id for update")
		return
	}
	post.ID = id
	if err := h.posts.update(post); err != nil {
		h.sendError(w, http.StatusInternalServerError, err, "error while update post")
		return
	}
	w.WriteHeader(http.StatusOK)
}

// api delete post
func (h *Handler) deletePostPage(w http.ResponseWriter, r *http.Request) {
	if err := h.posts.delete(chi.URLParam(r, "id")); err != nil {
		h.sendError(w, http.StatusInternalServerError, err, "error while delete post")
		return
	}
	w.WriteHeader(http.StatusOK)
}

// errors helper
func (h *Handler) sendError(w http.ResponseWriter, code int, err error, descr string) {
	log.Println(err, descr)
	w.WriteHeader(code)
	errMsg := Error{
		Code:  code,
		Err:   err.Error(),
		Descr: descr,
	}
	data, err := json.Marshal(errMsg)
	if err != nil {
		log.Println("Can't marshal error data:", err)
		return
	}
	if _, err = w.Write(data); err != nil {
		log.Println("Can't write to ResponseWriter:", err)
	}
}

// id counter
func (h *Handler) nextGlobID() int {
	h.mux.Lock()
	h.globID++
	h.mux.Unlock()
	return h.globID
}

// routes preparer
func (h *Handler) prepareRoutes() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.HandleFunc("/", h.mainPageForm)
	mux.Route(APIURL+POSTSURL, func(r chi.Router) {
		r.Post(CREATEURL, h.createPostPage)
		r.Put("/{id}", h.editPostPage)
		r.Delete("/{id}", h.deletePostPage)
	})
	mux.Route(POSTSURL, func(r chi.Router) {
		r.Get("/", h.postsPageForm)
		r.Get(CREATEURL, h.createPostPageForm)
		r.Get(EDITURL+"/", h.editPostPageForm)
	})
	return mux
}

// fill post map
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
