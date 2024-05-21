package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	"log"

	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UrlMappingDocument struct {
	ShortenUrl  string
	CompleteUrl string
}

type Form struct {
	URL string `form:"url"`
}

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

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"shortenUrl": nil,
		})
	})
	r.POST("/", func(c *gin.Context) {
		var form Form

		c.Bind(&form)

		if form.URL == "" {
			c.Status(400)
			return
		}

		mapping := UrlMappingDocument{
			ShortenUrl:  RandomId(6),
			CompleteUrl: form.URL,
		}

		_, err := collection.InsertOne(context.TODO(), mapping)

		if err != nil {
			c.Status(500)
			return
		}

		c.HTML(http.StatusCreated, "index.tmpl.html", gin.H{
			"shortenUrl": BuildShortenUrlFromId(mapping.ShortenUrl),
		})
	})
	r.GET("/r/:id", func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			c.Status(400)
		}

		var mapping UrlMappingDocument
		err := collection.FindOne(context.TODO(), bson.D{{"shortenurl", id}}).Decode(&mapping)

		if err == mongo.ErrNoDocuments {
			c.Status(http.StatusNotFound)
			return
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, mapping.CompleteUrl)

	})
	log.Fatal(r.Run())
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

func BuildShortenUrlFromId(id string) string {
	return "/r/" + id
}
