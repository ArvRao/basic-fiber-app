package main

import (
	"fmt"
	"os"

	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Book struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Title  string             `bson:"title,omitempty"`
	Author string             `bson:"author,omitempty"`
	Tags   []string           `bson:"tags,omitempty"`
	Pages  int                `bson:"pages,omitempty"`
}

func main() {
	app := fiber.New()
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://Arvind:mongoConnection@go-cluster.sbelezr.mongodb.net/go-db?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	database := client.Database("go-db")
	bookCollection := database.Collection("books")

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, Railway!",
		})
	})

	app.Get("/getAllNames", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"names": "Arvind",
		})
	})
	books := Book{
		Title:  "The Polyglot Developer",
		Author: "Nic Raboy",
		Tags:   []string{"development", "programming", "coding"},
		Pages:  278,
	}
	// fmt.Println("List of books: ", books)
	insertResult, err := bookCollection.InsertOne(ctx, books)
	if err != nil {
		panic(err)
	}
	fmt.Println("List of books: ", insertResult.InsertedID)

	app.Listen(getPort())
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	} else {
		port = ":" + port
	}

	return port
}
