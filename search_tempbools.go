package main

import (
	"log"
    "os"
	"bufio"
	"regexp"
	"sort"
	"strings"
	
	"github.com/gocarina/gocsv"
	
	"context"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Book struct {
	Id			int64	 `csv:"id"`
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
	Location   	string   `csv:"location"`
}

type BookValue struct {
	Book	map[string]interface {}
	value	int
}

type BookValues []BookValue
func (u BookValues) Len() int {
	return len(u)
}

func (u BookValues) Less(i, j int) bool {
	return u[i].value > u[j].value
}

func (u BookValues) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}


func main() {
	SearchBooks()
}

func SearchBooks() {
	offset := 1
	limit := 5

	// get DB data
	client, ctx := ConnectDB()
	defer client.Close()
	collection := "tmp_book"

	var bookvalues BookValues
	iter := client.Collection(collection).Where("exist", "==", "〇").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("【Error】", err)
			panic(err)
		}
		var book BookValue
		doc.DataTo(&book.Book)
		log.Printf("Book: %#v\n", book)
		bookvalues = append(bookvalues, book)
	}
	iter = client.Collection(collection).Where("exist", "==", "一部発見").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("【Error】", err)
			panic(err)
		}
		var book BookValue
		doc.DataTo(&book.Book)
		log.Printf("Book: %#v\n", book)
		bookvalues = append(bookvalues, book)
	}

	// get keyword list 
	stdin := bufio.NewScanner(os.Stdin)
	log.Printf("Please type keywords > ")
	stdin.Scan()
	splits := regexp.MustCompile("[ 　]+").Split(stdin.Text(), -1)
	text_split := []string{}
	for i := range splits {
		text_split = append(text_split, splits[i])
	}
	log.Printf("%d: %v", len(text_split), text_split)



	// search
	searchAttribute := []string{"publisher", "author", "bookName", "pubdate", "ISBN"}
	for i, book := range bookvalues {
		log.Printf("%d: %v", i, book)
		for _, word := range text_split {
			log.Printf("word: %s", word)
			for v, att := range searchAttribute {
				log.Printf("\tatt: %s", att)
				if strings.Index(book.Book[att].(string), word) != -1 {
					bookvalues[i].value += (v+1)
					break
				}
			}
		}
	}

	for _, book := range bookvalues {
		log.Printf("%v\n", book)
	}
	sort.Sort(bookvalues)
	log.Println("\nSorted")
	for _, book := range bookvalues {
		log.Printf("%v\n", book)
	}
	
	log.Println("\nRecommend")
	var books []map[string]interface {}
	for _, book := range bookvalues {
		if book.value != 0 {
			books = append(books, book.Book)
		}
	}
	for _, book := range books {
		log.Printf("%v\n", book)
	}

	log.Printf("\nReturn offset:%d, limit:%d, len:&d\n", offset, limit, len(books))
	if len(books) <= offset {
		log.Println("No Book!")
	} else {
		first := offset
		var last int
		if len(books) < offset+limit {
			last = len(books)
		}else{
			last = offset+limit
		}
		log.Printf("first:%d, last:%d\n", first, last)
		for _, book := range books[first:last] {
			log.Printf("%v\n", book)
		}
	}
}

func CreateTmpBooks() {
    file, err := os.OpenFile("../../nagaoBookList_format.csv", os.O_RDONLY, os.ModePerm)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    books := []*Book{}
    if err := gocsv.UnmarshalFile(file, &books); err != nil {
        panic(err)
	}
	log.Println(len(books))
	
    book_fires := []*BookFireStore{}
    for i, book := range books {
		var fire *BookFireStore = new(BookFireStore)
		fire.Id = book.Id
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
		fire.Location = "unidentified"
		book_fires = append(book_fires, fire)
		if i > 5 {
			break
		}
    }

	projectID := os.Getenv("LIBAPP_PROJECT_ID")
	collection := "tmp_book"
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	col := client.Collection(collection)
	log.Println("write")
	for i := range book_fires {
		log.Println(i, book_fires[i])
		_, _, err = col.Add(ctx, book_fires[i])
		if err != nil {
			log.Fatalf("Failed adding alovelace: %v", err)
		}
	}

	log.Println("read")
	iter := client.Collection(collection).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		log.Println(doc.Data())
	}
}
