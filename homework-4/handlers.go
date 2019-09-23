/*
 * HomeWork-4: Simple blog - MySQL
 * Created on 22.09.2019 13:11
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"log"
	"net/http"
	"sync"
)

// Handler is the global server handlers struct.
type Handler struct {
	db       *sql.DB
	posts    dbPosts
	tmplGlob *template.Template
	//globID   int
	mux sync.Mutex
	Error
}

// get all posts from DB
func (h *Handler) getAllPosts() error {
	rows, err := h.db.Query(GETALLPOSTS)
	if err != nil {
		return fmt.Errorf("error in all db.query: %v", err)
	}
	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Summary, &post.Body, &post.Date)
		if err != nil {
			return fmt.Errorf("error in all row.scan: %v", err)
		}
		h.posts = append(h.posts, post)
	}
	if err := rows.Close(); err != nil {
		return err
	}
	return nil
}

// get one post from DB
func (h *Handler) getOnePost(id string) (Post, error) {
	post := Post{}
	rows, err := h.db.Query(GETONEPOST, id)
	if err != nil {
		return post, fmt.Errorf("error in one db.query: %v", err)
	}
	rows.Next()
	if err := rows.Scan(&post.ID, &post.Title, &post.Summary, &post.Body, &post.Date); err != nil {
		return post, fmt.Errorf("error in one row.scan: %v", err)
	}
	if err := rows.Close(); err != nil {
		return post, fmt.Errorf("error in one row.close: %v", err)
	}
	return post, nil
}

// main page
func (h *Handler) mainPageForm(w http.ResponseWriter, _ *http.Request) {
	h.posts = dbPosts{}
	if err := h.getAllPosts(); err != nil {
		log.Println(err)
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
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
	post, err := h.getOnePost(postNum)
	if err != nil {
		log.Println(err)
		//http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	post.Body = template.HTML(blackfriday.Run([]byte(post.Body)))
	h.execTemplate(w, post, "post")
}

// crete post page
func (h *Handler) createPostPageForm(w http.ResponseWriter, _ *http.Request) {
	h.execTemplate(w, Post{}, "create")
}

// edit post page
func (h *Handler) editPostPageForm(w http.ResponseWriter, r *http.Request) {
	//post := h.posts[r.URL.Query().Get("id")]
	//h.execTemplate(w, post, "edit")
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
	//post := h.decodePost(w, r)
	//if post == nil {
	//	return
	//}
	//post.ID = h.nextGlobID()
	//if err := h.posts.create(post); err != nil {
	//	//if err := post.create(); err != nil {
	//	h.sendError(w, http.StatusInternalServerError, err, "error while create post")
	//	return
	//}
	//w.WriteHeader(http.StatusCreated)
}

// api edit post
func (h *Handler) editPostPage(w http.ResponseWriter, r *http.Request) {
	postNum := chi.URLParam(r, "id")
	post := h.decodePost(w, r)
	if post == nil {
		return
	}
	post.ID = postNum
	if err := h.posts.update(post); err != nil {
		//if err := post.update(); err != nil {
		h.sendError(w, http.StatusInternalServerError, err, "error while update post")
		return
	}
	w.WriteHeader(http.StatusOK)
}

// JSON decoder helper
func (h *Handler) decodePost(w http.ResponseWriter, r *http.Request) *Post {
	post := &Post{}
	err := json.NewDecoder(r.Body).Decode(post)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err, "error while decoding post body")
		return nil
	}
	return post
}

// api delete post
func (h *Handler) deletePostPage(w http.ResponseWriter, r *http.Request) {
	if err := h.posts.delete(chi.URLParam(r, "id")); err != nil {
		h.sendError(w, http.StatusInternalServerError, err, "error while delete post")
		return
	}
	w.WriteHeader(http.StatusOK)
}

// id counter
//func (h *Handler) nextGlobID() string {
//	h.mux.Lock()
//	h.globID++
//	h.mux.Unlock()
//	return strconv.Itoa(h.globID)
//}

// fill post map
//func (h *Handler) initPosts() {
//	h.posts = dbPosts{
//		"1": {
//			ID:      "1",
//			Title:   "Мой первый пост!",
//			Date:    "18-е Сентября 2019 года",
//			Summary: "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Odio praesentium, quos. Aspernatur assumenda cupiditate deserunt ducimus, eveniet, expedita inventore laboriosam magni modi non odio, officia qui sequi similique unde voluptatem.",
//			Body: `
//Здесь основной текст.
//# Markdown!
//*Это* **круто**!
//`,
//		},
//		"2": {
//			ID:    "2",
//			Title: "Это уже второй пост!",
//			Date:  "19-е Сентября 2019 года",
//			Summary: `
//Lorem ipsum dolor sit amet, consectetur adipisicing elit. Odio praesentium, quos. Aspernatur assumenda cupiditate deserunt ducimus, eveniet, expedita inventore laboriosam magni modi non odio, officia qui sequi similique unde voluptatem.
//Lorem ipsum dolor sit amet, consectetur adipisicing elit. Odio praesentium, quos. Aspernatur assumenda cupiditate deserunt ducimus, eveniet, expedita inventore laboriosam magni modi non odio, officia qui sequi similique unde voluptatem.
//`,
//			Body: `
//Разобрался в шаблонах и маркдаунах, как их совместить.
//
//Теперь понять, как переходить на отдельные посты.
//# Anybody!
//*Hz* **cool**!
//`,
//		},
//		"3": {
//			ID:      "3",
//			Title:   "Пример основных вариантов разметки Markdown",
//			Date:    "20-е Сентября 2019 года",
//			Summary: "Официальное руководство по синтаксису Markdown мне кажется слишком длинным и не слишком наглядным, поэтому я составил краткое руководство, которое поможет выучить или повторить синтаксис Маркдауна за час.",
//			Body:    `todo`,
//		},
//	}
//	h.globID = 3 // last id
//
//	items := []*Post{}
//
//	rows, err := h.DB.Query("SELECT id, title, updated FROM items")
//	__err_panic(err)
//	for rows.Next() {
//		post := &Item{}
//		err = rows.Scan(&post.Id, &post.Title, &post.Updated)
//		__err_panic(err)
//		items = append(items, post)
//	}
//	// надо закрывать соединение, иначе будет течь
//	rows.Close()
//}
