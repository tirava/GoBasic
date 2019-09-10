/*
 * HomeWork-2: Search string - 2: JSON server
 * Created on 15.09.19 12:12
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

var urls = [...]string{
	"https://www.litmir.me/br/?b=110008&p=1",
	"https://librebook.me/belyi_bim_chernoe_uho",
	"http://qqq.ww", // error here
	"https://knizhnik.org/dmitrij-gluhovskij/metro-2033/1",
	"https://www.gazeta.ru",
	"https://www.yandex.ru",
	"https://www.3dnews.ru",
}

func main() {

	//search := "Бим"
	//search := "Книга"
	//search := "книга"
	//search := "1973"
	//search := "2033"
	//search := "bug"

	//fmt.Println(searchStringURL(search, urls))

	data := []byte(`{"search":"bug"}`)
	r := bytes.NewReader(data)

	go runServer(":8080")

	resp, err := http.Post("http://localhost:8080", "application/json", r)
	if err != nil {
		log.Fatalln("error getting response from server:", err)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("error getting response body:", err)
	}

	fmt.Println(string(content))
}

func runServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				fmt.Fprint(w, "search response:", r.Body)
			}
		})

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Println("starting server at:", addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln("error starting server", err)
	}
}

func searchStringURL(search string) (res []string) {

	wg := &sync.WaitGroup{}
	mux := &sync.Mutex{}

	for _, url := range urls {

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
