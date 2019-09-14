/*
 * HomeWork-2: Routers
 * Created on 15.09.2019 22:02
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	addr        = "localhost:8080"
	cookieName  = "KlimGo"
	cookieValue = "GeekBrains"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {

	io.WriteString(w, "/write/...  - write cookie\n/read/...   - read cookie\n/delete/... - delete cookie")
}

func writeCookie(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:    cookieName,
		Value:   cookieValue,
		Expires: time.Now().AddDate(0, 0, 1),
		Path:    "/",
	})

	io.WriteString(w, "cookie written!")
}

func readCookie(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie(cookieName)
	if err != nil {
		fmt.Fprintln(w, "error reading cookie:", err)
		return
	}

	fmt.Fprintln(w, "read cookie:", c.Name, "=", c.Value)
}

func deleteCookie(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:   cookieName,
		Path:   "/",
		MaxAge: -1,
	})

	io.WriteString(w, "cookie deleted!")
}

func main() {

	// prepare server
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/write/", writeCookie)
	mux.HandleFunc("/read/", readCookie)
	mux.HandleFunc("/delete/", deleteCookie)

	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// safe shutdown
	go func() {
		<-shutdown
		// any work here
		fmt.Printf("\nShutdown server at: %s\n", addr)
		os.Exit(0)
	}()

	// start
	fmt.Println("Starting server at:", addr)
	log.Fatalln(http.ListenAndServe(addr, mux))
}
