/*
 * HomeWork-2: Search string - 2: JSON server
 * Created on 15.09.19 12:12
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
)

const (
	addr      = "localhost:8080"
	sitesFile = "sites.txt"
)

type search struct {
	Search string   `json:"search"`
	Sites  []string `json:"sites"`
}

var urls []string // contains sites urls after read from file

func handler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.Write([]byte("Hello and GoodBye! - Need POST method.\n"))
		return
	}

	// decode POST data
	s := &search{}
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		w.Write([]byte("Can't parse POST data.\n"))
		return
	}

	// get search results and code to json
	s.Sites = searchStringURL(s.Search, urls)
	b, err := json.MarshalIndent(s, "", "    ") // for best view in curl
	//b, err := json.Marshal(s)
	if err != nil {
		w.Write([]byte("Can't encode result data.\n"))
		return
	}

	// set proper headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(b)
}

func searchStringURL(search string, urls []string) (res []string) {

	wg := &sync.WaitGroup{}
	mux := &sync.Mutex{}

	for _, url := range urls {
		if len(url) < 3 { // no fake strings
			continue
		}

		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Error getting url: %v", err)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading body: %v", err)
				return
			}

			if strings.Contains(string(body), search) {
				mux.Lock()
				res = append(res, url)
				mux.Unlock()
			}
		}(url)
	}

	wg.Wait()
	return
}

func main() {

	// get URLs from file
	file, err := os.Open(sitesFile)
	if err != nil {
		log.Fatalln("Can't open file with sites:", sitesFile, err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln("Error reading site file body:", file, err)
	}

	urls = strings.Split(string(b), "\n")

	// prepare server
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, os.Kill)

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

// curl --header "Content-Type: application/json" --request POST --data '{"search":"bug"}' http://localhost:8080
// "Бим", "Книга", "книга", "1973", "2033"
