package handler

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"

	"go-url-shortener/types"
)

type CreateShortURLHandler struct {
	Collection       *mongo.Collection
	HomepageTemplate *template.Template
}

func (handler *CreateShortURLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlToShorten := r.FormValue("url")

	if urlToShorten == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "You must provide the 'url' form value")
		return
	}

	urlMapping := types.UrlMappingDocument{
		ShortenUrl:  RandomId(5),
		CompleteUrl: urlToShorten,
	}

	_, err := handler.Collection.InsertOne(r.Context(), urlMapping)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Cannot insert in the DB")
		return
	}

	w.Header().Add("Location", fmt.Sprintf("/shorten/%s", urlMapping.ShortenUrl))
	w.WriteHeader(http.StatusFound)
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
