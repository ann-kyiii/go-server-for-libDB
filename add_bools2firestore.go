package main

import (
	"log"
	"fmt"

    "os"
	"github.com/gocarina/gocsv"
	
	"context"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Book struct {
	BookName   	string   `csv:"bookName"`
	Genre      	string   `csv:"genre"`
	SubGenre   	string   `csv:"subGenre"`
	ISBN    	string   `csv:"ISBN"`
	Find 		int64    `csv:"find"`
	Sum 		int64    `csv:"sum"`
	Author    	string   `csv:"author"`
	Publisher   string   `csv:"publisher"`
	Pubdate 	string   `csv:"pubdate"`
	Exist 		string   `csv:"exist"`
	LocateAt4F  bool     `csv:"locateAt4F"`
	WithDisc    string   `csv:"withDisc"`
	Other    	string   `csv:"other"`
}
type Book_FireData struct {
	BookName   	string   `firestore:"bookName,omitempty"`
	Genre      	string   `firestore:"genre,omitempty"`
	SubGenre   	string   `firestore:"subGenre,omitempty"`
	ISBN    	string   `firestore:"ISBN,omitempty"`
	Find 		int64    `firestore:"find,omitempty"`
	Sum 		int64    `firestore:"sum,omitempty"`
	Author    	string   `firestore:"author,omitempty"`
	Publisher   string   `firestore:"publisher,omitempty"`
	Pubdate		string   `firestore:"pubdate,omitempty"`
	Exist 		string   `firestore:"exist,omitempty"`
	LocateAt4F  bool     `firestore:"locateAt4F,omitempty"`
	WithDisc    string   `firestore:"withDisc,omitempty"`
	Other    	string   `firestore:"other,omitempty"`
	Borrower	interface{} `firestore:"borrower,omitempty"`
}


func main() {
    file, err := os.OpenFile("../nagaoBookList_format.csv", os.O_RDONLY, os.ModePerm)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    books := []*Book{}
    if err := gocsv.UnmarshalFile(file, &books); err != nil {
        panic(err)
	}
	fmt.Println(len(books))
	
    book_fires := []*Book_FireData{}
    for _, book := range books {
		var fire *Book_FireData = new(Book_FireData)
		fire.BookName = book.BookName
		fire.Genre = book.Genre
		fire.SubGenre = book.SubGenre
		fire.ISBN = book.ISBN
		fire.Find = book.Find
		fire.Sum = book.Sum
		fire.Author = book.Author
		fire.Publisher = book.Publisher
		fire.Pubdate = book.Pubdate
		fire.Exist = book.Exist
		fire.LocateAt4F = book.LocateAt4F
		fire.WithDisc = book.WithDisc
		fire.Other = book.Other
		fire.Borrower = []string{}
		book_fires = append(book_fires, fire)
        // fmt.Println(book)
    }

	projectID := "nagaolab-library-app" // Sets your Google Cloud Platform project ID.
	// Get a Firestore client.
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close() // Close client when done.

	col := client.Collection("nagao_books")
	fmt.Println("write")
	for i := range book_fires {
		fmt.Println(i, book_fires[i])
		_, _, err = col.Add(ctx, book_fires[i])
		if err != nil {
			log.Fatalf("Failed adding alovelace: %v", err)
		}
		// if i > 0{
		// 	break 
		// }
	}
	fmt.Println("read")
	iter := client.Collection("nagao_books").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Println(doc.Data())
	}
}
