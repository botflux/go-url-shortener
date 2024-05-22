package main

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"strings"

	"log"

	"math/rand"
	"net/http"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatal("Set the 'MONGODB_URI' environment variable to start the application")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatalln("Cannot connect to MongoDB using the connection string passed by environment variable", err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalln("Something went wrong while disconnecting from MongoDB", err)
		}
	}()

	collection := client.Database("url_shortener").Collection("url_mapping")

	name, err := collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{Key: "shortenurl", Value: 1}},
		Options: options.Index().SetUnique(true).SetExpireAfterSeconds(60),
	})

	if err != nil {
		log.Fatalln("Cannot create indexes on the MongoDB collection", err)
	}

	fmt.Printf("Index successfully created with name '%s'\n", name)

	indexTemplate := template.Must(template.ParseFiles("templates/index.tmpl.html"))

	http.Handle("GET /", &HomepageHandler{
		HomepageTemplate: indexTemplate,
	})

	http.Handle("POST /", &CreateShortURLHandler{
		Collection:       collection,
		HomepageTemplate: indexTemplate,
	})

	http.Handle("GET /r/{id}", &RedirectToCompleteURLHandler{
		Collection: collection,
	})

	server := &http.Server{
		Addr: ":4500",
	}

	log.Fatal(server.ListenAndServe())
}

var letters = [...]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

func RandomId(length int) string {
	var sb strings.Builder

	for i := 0; i < length; i++ {
		index := rand.Intn(25)

		sb.WriteString(
			letters[index],
		)
	}

	return sb.String()
}
