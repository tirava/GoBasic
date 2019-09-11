/*
 * HomeWork-1: Search string
 * Created on 11.09.19 22:41
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

func main() {

	urls := []string{
		"https://www.litmir.me/br/?b=110008&p=1",
		"https://librebook.me/belyi_bim_chernoe_uho",
		"http://qqq.ww", // error here
		"https://knizhnik.org/dmitrij-gluhovskij/metro-2033/1",
		"https://www.gazeta.ru",
		"https://www.yandex.ru",
		"https://www.3dnews.ru",
	}

	//search := "Бим"
	//search := "Книга"
	//search := "книга"
	//search := "1973"
	//search := "2033"
	search := "bug"

	fmt.Printf("found string '%s' in sites: %v\n", search, searchStringURL(search, urls))
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
