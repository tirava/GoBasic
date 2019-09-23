/*
 * HomeWork-4: Simple blog - MySQL
 * Created on 22.09.2019 13:11
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
)

// Constants
const (
	SERVADDR     = ":8080"
	TEMPLATEEXT  = "*.gohtml"
	TEMPLATEPATH = "templates"
	POSTSURL     = "/posts"
	EDITURL      = "/edit"
	CREATEURL    = "/create"
	APIURL       = "/api/v1"
	STATICPATH   = "/static"
)

func main() {

	// new handlers struct
	handlers := &Handler{
		tmplGlob: template.Must(template.ParseGlob(path.Join(TEMPLATEPATH, TEMPLATEEXT))),
	}

	// fill posts
	handlers.initPosts()

	// prepare server, routes & middleware
	srv := &http.Server{Addr: SERVADDR, Handler: handlers.prepareRoutes()}

	// graceful shutdown
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM) // os.Kill cannot be trapped anyway!
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		log.Println("Signal received:", <-shutdown)
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("Error while shutdown server:", err)
		}
	}()

	fmt.Println("Starting server at:", SERVADDR)
	log.Printf("Shutdown server at: %s\n%v", SERVADDR, srv.ListenAndServe())
}
