/*
 * HomeWork-2: Search string - 2: JSON server
 * Created on 15.09.19 12:12
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

const addr = "localhost:8080"

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Write([]byte("Hello and GoodBye!"))
		return
	}
	fmt.Fprint(w, "search response:", r.Body)
}

func searchStringURL(search string, urls []string) (res []string) {

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

func main() {

	/*var urls = []string{
		"https://www.litmir.me/br/?b=110008&p=1",
		"https://librebook.me/belyi_bim_chernoe_uho",
		"http://qqq.ww", // error here
		"https://knizhnik.org/dmitrij-gluhovskij/metro-2033/1",
		"https://www.gazeta.ru",
		"https://www.yandex.ru",
		"https://www.3dnews.ru",
	}*/

	//search := "Бим"
	//search := "Книга"
	//search := "книга"
	//search := "1973"
	//search := "2033"
	//search := "bug"

	//fmt.Println(searchStringURL(search, urls))

	//data := []byte(`{"search":"bug"}`)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	fmt.Println("Starting server at:", addr)
	log.Fatalln(http.ListenAndServe(addr, mux))

}
