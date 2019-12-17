package main

import (
	"os"
	"strconv"
	"log"
	"fmt"
	"context"
	
	"net/http"
    "github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"

	"google.golang.org/api/iterator"
	"github.com/ipfans/echo-session"
	"cloud.google.com/go/firestore"
)



func main() {
	e := echo.New()

    store := session.NewCookieStore([]byte("secret-key"))
    store.MaxAge(2)
	e.Use(session.Sessions("ESESSION", store))
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
	e.GET("/", hello, ContinuousAccessFilter())
	e.GET("/api/v1/bookId/:id", getBookWithID, ContinuousAccessFilter())
	e.POST("/api/v1/search", searchBooks, ContinuousAccessFilter())
	e.POST("/api/v1/searchGenre", searchGenre, ContinuousAccessFilter())
	e.POST("/api/v1/searchSubGenre", searchSubGenre, ContinuousAccessFilter())
	e.POST("/api/v1/borrow", borrowBook, ContinuousAccessFilter())
	e.POST("/api/v1/return", returnBook, ContinuousAccessFilter())
}

func ContinuousAccessFilter() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
			session := session.Default(c)
			access := session.Get("AccessServer")
			if access != nil && access == "completed" {
				return echo.ErrUnauthorized
			}
			
			return next(c)
        }
    }
}

func hello(c echo.Context) error {
    session := session.Default(c)
    session.Set("AccessServer", "completed")
    session.Save()
	return c.JSON(http.StatusOK, map[string]string{"hello": "world"})
}


func getBookWithID(c echo.Context) error {
    session := session.Default(c)
    session.Set("AccessServer", "completed")
    session.Save()
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
	var book map[string]interface {}
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
    session := session.Default(c)
    session.Set("AccessServer", "completed")
    session.Save()
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
    session := session.Default(c)
    session.Set("AccessServer", "completed")
    session.Save()
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

	var books []map[string]interface {}
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
		var book map[string]interface {}
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
		var book map[string]interface {}
		doc.DataTo(&book)
		books = append(books, book)
	}

	if len(books) <= offset {
		empty_list := []interface{}{}
		data := map[string]interface{}{
			"books":empty_list,
			"max_books":0,
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
			"max_books":len(books),
		}
		return c.JSON(http.StatusOK, data)
	}
}

func searchSubGenre(c echo.Context) error {
    session := session.Default(c)
    session.Set("AccessServer", "completed")
    session.Save()
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

	var books []map[string]interface {}
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
		var book map[string]interface {}
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
		var book map[string]interface {}
		doc.DataTo(&book)
		books = append(books, book)
	}

	if len(books) <= offset {
		empty_list := []interface{}{}
		data := map[string]interface{}{
			"books":empty_list,
			"max_books":0,
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
			"max_books":len(books),
		}
		return c.JSON(http.StatusOK, data)
	}
}

func borrowBook(c echo.Context) error {
	session := session.Default(c)
    session.Set("AccessServer", "completed")
    session.Save()
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	t_id := m["id"].(string)
	name := m["name"].(string)
	id, err := strconv.Atoi(t_id)
	if err != nil {
		log.Printf("【Error】", err)
		panic(err)
	}

	// get DB data
	client, ctx := ConnectDB()
	defer client.Close()
	data := borrow_book(ctx, client, id, name)

	return c.JSON(http.StatusOK, data)
}

func returnBook(c echo.Context) error {
	session := session.Default(c)
    session.Set("AccessServer", "completed")
    session.Save()
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	t_id := m["id"].(string)
	name := m["name"].(string)
	id, err := strconv.Atoi(t_id)
	if err != nil {
		log.Printf("【Error】", err)
		panic(err)
	}
	
	// get DB data
	client, ctx := ConnectDB()
	defer client.Close()
	data := return_book(ctx, client, id, name)

	return c.JSON(http.StatusOK, data)
}

func borrow_book(ctx context.Context, client *firestore.Client, id int, name string) map[string]interface{}{
	collection := os.Getenv("LIBAPP_COLLECTION")
	var book map[string]interface {}

	doc := client.Collection(collection).Doc(strconv.Itoa(id))

	docu, _ := doc.Get(ctx)
	d := docu.Data()
	borrower := d["borrower"].([]interface{})
	arr := make([]string, len(borrower))
	for i, v := range borrower{
		arr[i] = fmt.Sprint(v)
	}
	arr = append(arr, name)

	_, err := doc.Set(ctx, map[string]interface{}{
		"borrower": arr,
	}, firestore.MergeAll)
	if err != nil{
		log.Printf("An error has occurred: %s", err)
	}
	iter := client.Collection(collection).Where("id", "==", id).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		doc.DataTo(&book)
	}
	return book
}

func remove(strings []string, search string) []string {
    result := []string{}
    for _, v := range strings {
        if v != search {
            result = append(result, v)
        }
    }
    return result
}

func return_book(ctx context.Context, client *firestore.Client, id int, name string) map[string]interface{}{
	collection := os.Getenv("LIBAPP_COLLECTION")
	var book map[string]interface {}
	doc := client.Collection(collection).Doc(strconv.Itoa(id))

	docu, _ := doc.Get(ctx)
	d := docu.Data()
	borrower := d["borrower"].([]interface{})
	arr := make([]string, len(borrower))
	for i, v := range borrower{
		arr[i] = fmt.Sprint(v)
	}
	arr = remove(arr, name)

	_, err := doc.Set(ctx, map[string]interface{}{
		"borrower": arr,
	}, firestore.MergeAll)
	if err != nil{
		log.Printf("An error has occurred: %s", err)
	}
	iter := client.Collection(collection).Where("id", "==", id).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		doc.DataTo(&book)
	}
	return book
}
