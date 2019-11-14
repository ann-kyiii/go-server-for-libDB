package main

import (
	"log"
	"fmt"
	"os"
	"bufio"
	"strconv"
	// "reflect"

	"context"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func borrow(ctx context.Context, client *firestore.Client, id int) error {
	var sc = bufio.NewScanner(os.Stdin)
	var name string
	fmt.Printf("Please type Borrower Name > ")
	if sc.Scan(){
		name = sc.Text()
	}
	fmt.Println("BORROWER : " + name)
	doc := client.Collection("tmp_book").Doc(strconv.Itoa(id))

	docu, _ := doc.Get(ctx)
	d := docu.Data()
	borrower := d["borrower"].([]interface{})
	// sum := d["sum"]
	arr := make([]string, len(borrower))
	for i, v := range borrower{
		arr[i] = fmt.Sprint(v)
	}
	arr = append(arr, name)

	// fmt.Println(arr)
	// fmt.Println(reflect.TypeOf(arr))

	_, err := doc.Set(ctx, map[string]interface{}{
		"borrower": arr,
	}, firestore.MergeAll)
	if err != nil{
		log.Printf("An error has occurred: %s", err)
	}
	return err
}

func main() {
	projectID := os.Getenv("LIBAPP_PROJECT_ID") // Sets your Google Cloud Platform project ID.
	// collection := os.Getenv("LIBAPP_COLLECTION")
	collection := "tmp_book"
	// Get a Firestore client.
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close() // Close client when done.

	col := client.Collection(collection)
	// stdin := bufio.NewScanner(os.Stdin)

	var sc = bufio.NewScanner(os.Stdin)
	var text string
	fmt.Printf("Please type BookID > ")
	if sc.Scan(){
		text = sc.Text()
	}
	id, err := strconv.Atoi(text)
	borrow(ctx, client, id)
	iter := col.Where("id", "==", id).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Printf("Please type Borrower Name > ")
			
		var book BookFireStore
		doc.DataTo(&book)
		fmt.Printf("Book: %#v\n", book)
	}
}
