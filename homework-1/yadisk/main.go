/*
 * HomeWork-1: Yandex file
 * Created on 13.09.19 23:15
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	fileURL := "https://yadi.sk/i/pBfU5WBqFWO0FA"

	fileName, err := getFileFromURL(fileURL)
	if err != nil {
		log.Fatalf("Error while download file from URL:\n%s \n%v", fileURL, err)
	}

	fmt.Println(fileName)
}

func getFileFromURL(url string) (string, error) {
	const reqURL = "https://cloud-api.yandex.net/v1/disk/public/resources/download?public_key="

	resp, err := http.Get(reqURL + url)
	if err != nil {
		return "", errors.New("error getting file")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("file not found")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("error reading file body")
	}

	return string(body), nil
}
