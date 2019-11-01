package main

import (
	"os"
	
	"context"
	"cloud.google.com/go/firestore"
)


func ConnectDB() (*firestore.Client, context.Context)  {
	projectID := os.Getenv("LIBAPP_PROJECT_ID") // Sets your Google Cloud Platform project ID.
	ctx := context.Background() // Get a Firestore client.
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
		// log.Fatalf("Failed to create client: %v", err)
	}

	return client, ctx
}
