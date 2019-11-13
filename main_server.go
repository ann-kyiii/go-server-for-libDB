package main

import (
	"os"
	"strconv"
	"log"
	
	"net/http"
    "github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"

	"google.golang.org/api/iterator"
)



func main() {
	e := echo.New()

    e.Use(middleware.CORS())
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		log.Printf("Request: %v\n", string(reqBody))
	}))
	initRouting(e)

	e.Logger.Fatal(e.Start(":1313"))
}

func initRouting(e *echo.Echo) {
	e.GET("/", hello)
	e.GET("/api/v1/bookId/:id", getBookWithID)
	e.POST("/api/v1/search", searchBooks)
	e.POST("/api/v1/searchGenre", searchGenre)
	e.POST("/api/v1/searchSubGenre", searchSubGenre)
}


func hello(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"hello": "world"})
}


func getBookWithID(c echo.Context) error {
	id := c.Param("id")

	bookId, err := strconv.Atoi(id)
	if err != nil {
		return errors.Wrapf(err, "errors when book id convert to int: %s", bookId)
	}
	
	client, ctx := ConnectDB()
	defer client.Close() // Close client when done.
	collection := os.Getenv("LIBAPP_COLLECTION")
	col := client.Collection(collection)
	iter := col.Where("id", "==", bookId).Documents(ctx)
	var book BookFireStore
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("【Error】", err)
			return errors.Wrapf(err, "errors when getting bookInfo from DB")
		}
		doc.DataTo(&book)
		log.Printf("Book: %#v\n", book)
	}
	return c.JSON(http.StatusOK, book)
}


func searchBooks(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	keywords := m["keywords"].([]interface{})
	t_offset := m["offset"].(string)
	t_limit := m["limit"].(string)
	offset, err := strconv.Atoi(t_offset)
	if err != nil {
		log.Printf("【Error】", err)
		panic(err)
	}
	limit, err2 := strconv.Atoi(t_limit)
	if err2 != nil {
		log.Printf("【Error】", err)
		panic(err)
	}

	// get DB data
	client, ctx := ConnectDB()
	defer client.Close()
	collection := os.Getenv("LIBAPP_COLLECTION")

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
		bookvalues = append(bookvalues, book)
	}

	// search
	searchAttribute := []string{"publisher", "author", "bookName", "pubdate", "ISBN"}
	data := searchOR(bookvalues, keywords, searchAttribute, offset, limit)
	return c.JSON(http.StatusOK, data)
}

func searchGenre(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	genre := m["genre"].(string)
	t_offset := m["offset"].(string)
	t_limit := m["limit"].(string)
	offset, err := strconv.Atoi(t_offset)
	if err != nil {
		log.Printf("【Error】", err)
		panic(err)
	}
	limit, err2 := strconv.Atoi(t_limit)
	if err2 != nil {
		log.Printf("【Error】", err)
		panic(err)
	}

	// get DB data
	client, ctx := ConnectDB()
	defer client.Close()
	collection := os.Getenv("LIBAPP_COLLECTION")

	var books []BookFireStore
	iter := client.Collection(collection).Where("exist", "==", "〇").Where("genre", "==", genre).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("【Error】", err)
			panic(err)
		}
		var book BookFireStore
		doc.DataTo(&book)
		books = append(books, book)
	}
	iter = client.Collection(collection).Where("exist", "==", "一部発見").Where("genre", "==", genre).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("【Error】", err)
			panic(err)
		}
		var book BookFireStore
		doc.DataTo(&book)
		books = append(books, book)
	}

	if len(books) <= offset {
		empty_list := []interface{}{}
		data := map[string]interface{}{
			"books":empty_list,
		}
		return c.JSON(http.StatusOK, data)
	} else {
		first := offset
		var last int
		if len(books) < offset+limit {
			last = len(books)
		}else{
			last = offset+limit
		}
		data := map[string]interface{}{
			"books":books[first:last],
		}
		return c.JSON(http.StatusOK, data)
	}
}

func searchSubGenre(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	subGenre := m["subGenre"].(string)
	t_offset := m["offset"].(string)
	t_limit := m["limit"].(string)
	offset, err := strconv.Atoi(t_offset)
	if err != nil {
		log.Printf("【Error】", err)
		panic(err)
	}
	limit, err2 := strconv.Atoi(t_limit)
	if err2 != nil {
		log.Printf("【Error】", err)
		panic(err)
	}

	// get DB data
	client, ctx := ConnectDB()
	defer client.Close()
	collection := os.Getenv("LIBAPP_COLLECTION")

	var books []BookFireStore
	iter := client.Collection(collection).Where("exist", "==", "〇").Where("subGenre", "==", subGenre).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("【Error】", err)
			panic(err)
		}
		var book BookFireStore
		doc.DataTo(&book)
		books = append(books, book)
	}
	iter = client.Collection(collection).Where("exist", "==", "一部発見").Where("subGenre", "==", subGenre).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("【Error】", err)
			panic(err)
		}
		var book BookFireStore
		doc.DataTo(&book)
		books = append(books, book)
	}

	if len(books) <= offset {
		empty_list := []interface{}{}
		data := map[string]interface{}{
			"books":empty_list,
		}
		return c.JSON(http.StatusOK, data)
	} else {
		first := offset
		var last int
		if len(books) < offset+limit {
			last = len(books)
		}else{
			last = offset+limit
		}
		data := map[string]interface{}{
			"books":books[first:last],
		}
		return c.JSON(http.StatusOK, data)
	}
}
