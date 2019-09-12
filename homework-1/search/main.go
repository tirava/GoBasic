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

	// get URLs from file
	//urls := make([]siteMirror, 0)
	//
	//file, err := os.Open(fileName)
	//check(err, "fatal", "Can't open file with sites!")
	//defer file.Close()
	//
	//// read sites by line
	//f := bufio.NewReader(file)
	//for {
	//	line, err := f.ReadString('\n')
	//	if err == io.EOF {
	//		break
	//	}
	//	if len(line) < 3 { // no fake symbols
	//		continue
	//	}
	//	sitesNames = append(sitesNames, siteMirror{strings.TrimRight(line, "\n"), 0})
	// may be readall and split \n?
	//	}

	urls := []string{
		"https://www.litmir.me/br/?b=110008&p=1",
		"https://librebook.me/belyi_bim_chernoe_uho",
		"http://qqq.ww", // error here
		"https://knizhnik.org/dmitrij-gluhovskij/metro-2033/1",
		"https://www.gazeta.ru",
		"https://www.yandex.ru",
		"https://www.3dnews.ru",
	}

	//examples: "Бим", "Книга", "1973", "2033", "bug"

	search := ""
	for {
		fmt.Printf("Enter search URL (Ctrl-C for exit): ")
		_, err := fmt.Scanln(&search)
		if err != nil {
			log.Println("error parse search string", err)
			continue
		}

		fmt.Printf("Found string '%s' in sites:\n", search)
		found := searchStringURL(search, urls)

		for _, f := range found {
			fmt.Println(f)
		}
	}
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
