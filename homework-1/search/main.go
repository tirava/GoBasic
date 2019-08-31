/*
 * HomeWork-1: Search string
 * Created on 30.08.19 22:41
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
		"http://qqq.ww",
		"https://knizhnik.org/dmitrij-gluhovskij/metro-2033/1",
	}

	search := "Бим"
	//search := "1973"
	//search := "2033"

	fmt.Println(searchStringURL(search, urls))
}

func searchStringURL(search string, urls []string) (res []string) {

	wg := &sync.WaitGroup{}

	for _, url := range urls {
		//resp, err := http.Get(url)
		//if err != nil {
		//	log.Printf("Error getting url: %v", err)
		//	continue
		//}
		//body, err := ioutil.ReadAll(resp.Body)
		//if err != nil {
		//	log.Printf("Error reading body: %v", err)
		//	continue
		//}
		//if strings.Contains(string(body), search) {
		//	res = append(res, url)
		//}
		//_ = resp.Body.Close()

		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Error getting url: %v", err)
				return
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading body: %v", err)
				return
			}
			if strings.Contains(string(body), search) {
				res = append(res, url)
			}
			_ = resp.Body.Close()
		}(url)
	}

	wg.Wait()
	return
}
