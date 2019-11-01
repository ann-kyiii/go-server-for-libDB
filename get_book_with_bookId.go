package main

import (
	"log"
	"fmt"
	"os"
	"bufio"
	"strconv"

	"context"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func main() {
	projectID := os.Getenv("LIBAPP_PROJECT_ID") // Sets your Google Cloud Platform project ID.
	collection := os.Getenv("LIBAPP_COLLECTION")
	// Get a Firestore client.
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close() // Close client when done.

	col := client.Collection(collection)
	stdin := bufio.NewScanner(os.Stdin)

	
	fmt.Printf("Please type BookID > ")
	for stdin.Scan(){
		text := stdin.Text()
		id, err := strconv.Atoi(text)
		if err != nil  {
			break
		}
		iter := col.Where("id", "==", id).Documents(ctx)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
			}
			
			var book BookFireStore
			doc.DataTo(&book)
			fmt.Printf("Book: %#v\n", book)
		}
		fmt.Printf("Please type BookID > ")
	}
}
