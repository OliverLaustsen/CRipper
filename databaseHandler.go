package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func CreateDatabaseClient(ctx context.Context) (*firestore.Client, *auth.Client) {
	conf := &firebase.Config{ProjectID: GetEnvVariable("PROJECT_ID")}
	sa := "./serviceAccount.json"
	auth := option.WithCredentialsFile(sa)
	app, err := firebase.NewApp(ctx, conf, auth)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	authClient, authErr := app.Auth(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	if authErr != nil {
		log.Fatalln(authErr)
	}

	return client, authClient
}

func UploadFile(ctx context.Context, fromFilePath string, toFilePath string) chan string {
	c := make(chan string, 1)
	bName := GetEnvVariable("STORAGE_BUCKET")
	config := &firebase.Config{
		StorageBucket: bName,
	}
	app, err := firebase.NewApp(context.Background(), config, option.WithCredentialsFile("./serviceAccount.json"))
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		log.Fatalf("error initializing storage client: %v\n", err)
	}

	file, err := os.Open(fromFilePath)
	if err != nil {
		log.Fatalf("error opening file: %v\n", err)
	}
	defer file.Close()

	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Fatalf("error getting default bucket: %v\n", err)
	}

	object := bucket.Object(toFilePath)
	writer := object.NewWriter(ctx)

	if _, err := io.Copy(writer, file); err != nil {
		log.Fatalf("error uploading file: %v\n", err)
	}

	if err := writer.Close(); err != nil {
		log.Fatalf("error closing writer: %v\n", err)
	}
	c <- "File uploaded successfully"

	return c
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
