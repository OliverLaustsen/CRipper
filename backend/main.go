package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	godotenv.Load()

	dbc, _ := CreateDatabaseClient(ctx)
	defer dbc.Close()

	sources := []string{"https://asura.gg/"}
	comicMap := createComicMap()
	visitList := ScrapeSiteForReleases(sources, comicMap)
	dbc.Collection("comics").NewDoc().Create(ctx, comicMap)
	var wg sync.WaitGroup

	for _, e := range visitList {
		wg.Add(1)
		e := e
		go func() {
			ScrapeImage(e, ctx)
			wg.Done()
		}()
	}
	wg.Wait()
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

func GetEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
