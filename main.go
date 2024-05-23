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
	"go.uber.org/zap"

	"go-url-shortener/handler"
)

func main() {
	logger, err := zap.NewProduction()

	if err != nil {
		log.Fatalln("Cannot create the logger", err)
	}

	defer logger.Sync()

	logger.Info("App is starting up!")

	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found")
	}

	client, err := CreateMongoClientFromEnvs()

	if err != nil {
		logger.Panic("Cannot create the MongoDB client", zap.Error(err))
	}

	defer DisconnectMongoClient(client)

	collection, err := GetCollectionAndCreateIndices(client, logger)

	if err != nil {
		logger.Panic("Cannot get the collection or create the indices", zap.Error(err))
	}

	indexTemplate := template.Must(template.ParseFiles("templates/index.tmpl.html"))

	http.Handle("GET /", &handler.HomepageHandler{
		HomepageTemplate: indexTemplate,
	})

	http.Handle("POST /", &handler.CreateShortURLHandler{
		Collection:       collection,
		HomepageTemplate: indexTemplate,
	})

	http.Handle("GET /shorten/{id}", &handler.DisplayShortenUrlHandler{
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

	if err := server.ListenAndServe(); err != nil {
		logger.Panic("HTTP server stopped", zap.Error(err))
	}
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

func GetCollectionAndCreateIndices(client *mongo.Client, logger *zap.Logger) (*mongo.Collection, error) {
	collection := client.Database("url_shortener").Collection("url_mapping")

	name, err := collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{Key: "shortenurl", Value: 1}},
		Options: options.Index().SetUnique(true).SetExpireAfterSeconds(60),
	})

	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("Index successfully created with name '%s'", name))

	return collection, nil
}
