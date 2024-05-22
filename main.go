package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"os"

	"log"

	"net/http"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-url-shortener/handler"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	client, err := CreateMongoClientFromEnvs()

	if err != nil {
		log.Fatalln("Cannot create the MongoDB client", err)
	}

	defer DisconnectMongoClient(client)

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

	http.Handle("GET /", &handler.HomepageHandler{
		HomepageTemplate: indexTemplate,
	})

	http.Handle("POST /", &handler.CreateShortURLHandler{
		Collection:       collection,
		HomepageTemplate: indexTemplate,
	})

	http.Handle("GET /r/{id}", &handler.RedirectToCompleteURLHandler{
		Collection: collection,
	})

	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	server := &http.Server{
		Addr: ":4500",
	}

	log.Fatal(server.ListenAndServe())
}

func CreateMongoClientFromEnvs() (*mongo.Client, error) {
	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		return nil, errors.New("Set the 'MONGODB_URI' environment variable to start the application")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		return nil, err
	}

	return client, nil
}

func DisconnectMongoClient(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatalln("Something went wrong while disconnecting from MongoDB", err)
	}
}
