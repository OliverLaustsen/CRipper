package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
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

func GetEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func UploadFile(w io.Writer, bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	// Open local file.
	f, err := os.Open("notes.txt")
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	o = o.If(storage.Conditions{DoesNotExist: true})
	// If the live object already exists in your bucket, set instead a
	// generation-match precondition using the live object's generation number.
	// attrs, err := o.Attrs(ctx)
	// if err != nil {
	//      return fmt.Errorf("object.Attrs: %w", err)
	// }
	// o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// Upload an object with storage.Writer.
	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}
	fmt.Fprintf(w, "Blob %v uploaded.\n", object)
	return nil
}
