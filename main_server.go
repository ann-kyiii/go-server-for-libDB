package main

import (
	"os"
	"strconv"
	"log"
	
	"net/http"
    "github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"

	"context"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)



func main() {
	e := echo.New()

    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
	
	initRouting(e)

	e.Logger.Fatal(e.Start(":1313"))
}

func initRouting(e *echo.Echo) {
	e.GET("/", hello)
	e.GET("/api/v1/bookId/:id", getBookWithID)
}

func ConnectDB() (*firestore.Client, context.Context)  {
	projectID := os.Getenv("LIBAPP_PROJECT_ID") // Sets your Google Cloud Platform project ID.
	log.Printf("DB: %s\n", projectID)
	ctx := context.Background() // Get a Firestore client.
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
		// log.Fatalf("Failed to create client: %v", err)
	}

	return client, ctx
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
	var book Book_FireStore
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
