package main

import (
	"context"
	"log"

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

// func CreateStorageClient(ctx context.Context) {
// 	bName := GetEnvVariable("STORAGE_BUCKET")

// 	client := storage.NewClient(ctx)
// }

// func CreateStorageClient(ctx context.Context) {

// 	opt := option.WithCredentialsFile("path/to/serviceAccountKey.json")
// 	app, err := firebase.NewApp(context.Background(), config, opt)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	client, err := app.Storage(context.Background())
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	bucket, err := client.Bucket(bName).ob
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	return bucket
// }
