package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {

	sources := []string{"https://asura.gg/"}
	comicMap := createComicMap()
	visitList := ScrapeSiteForReleases(sources, comicMap)

	c := make(chan string, 2)

	for _, e := range visitList {

		ScrapeImage(e, c)
	}
	fmt.Println("Complete")
}

func createComicMap() map[string]int64 {
	cm := make(map[string]int64)

	data, err := os.ReadFile("./comiclist.txt")
	if err != nil {
		fmt.Println("Error while reading file")
	}

	entries := strings.Split(string(data), "\n")
	for _, e := range entries {
		el := strings.Split(e, ":-:")
		pnum, err := strconv.ParseInt(el[1], 10, 0)
		if err != nil {
			fmt.Println("Error while parsing number", err)
		}
		cm[el[0]] = pnum
	}

	return cm
}

func DownloadFile(URL, fileName string) chan string {
	r := make(chan string, 1)
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Println("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println(err)
	}
	r <- "Download of " + URL + " is finished"
	return r
}
