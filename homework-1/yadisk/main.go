/*
 * HomeWork-1: Yandex file
 * Created on 13.09.19 23:15
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const reqURL = "https://cloud-api.yandex.net/v1/disk/public/resources/download?public_key="

// YandexFile is the base struct for downloading files by public URL.
type yandexFile struct {
	Href      string
	Method    string
	Templated bool
}

func main() {

	fileURL := "https://yadi.sk/i/pBfU5WBqFWO0FA"

	fileName, err := getFileFromURL(fileURL)
	if err != nil {
		log.Fatalf("Error while download file from URL:\n%s \n%v", fileURL, err)
	}

	fmt.Println(fileName)
}

func getFileFromURL(url string) (string, error) {

	// get file metadata
	resp, err := http.Get(reqURL + url)
	if err != nil {
		return "", errors.New("error getting file metadata")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("file not found")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("error reading file metadata body")
	}

	// unmarshal metadata - get direct link
	yf := &yandexFile{}

	err = json.Unmarshal([]byte(body), yf)
	if err != nil {
		return "", errors.New("error unmarshal json metadata")
	}

	// get file body
	respFile, err := http.Get(yf.Href)
	if err != nil {
		return "", errors.New("error getting file")
	}
	defer respFile.Body.Close()

	bodyFile, err := ioutil.ReadAll(respFile.Body)
	if err != nil {
		return "", errors.New("error reading file body")
	}

	// save file
	err = ioutil.WriteFile("111.odt", bodyFile, 0644)
	if err != nil {
		return "", errors.New("error writing file on disk")
	}

	return "", nil
}
