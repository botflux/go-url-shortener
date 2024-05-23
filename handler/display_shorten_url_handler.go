package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go-url-shortener/types"
)

type DisplayShortenUrlHandler struct {
	Collection       *mongo.Collection
	HomepageTemplate *template.Template
}

func (handler *DisplayShortenUrlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	var mapping types.UrlMappingDocument

	err := handler.Collection.FindOne(r.Context(), bson.D{{Key: "shortenurl", Value: id}}).Decode(&mapping)

	if err == mongo.ErrNoDocuments {
		fmt.Fprint(w, "No mapping found for the given ID")
	}

	if err != nil {
		fmt.Fprint(w, "Cannot retrieve the mapping")
	}

	w.WriteHeader(http.StatusOK)
	handler.HomepageTemplate.Execute(w, mapping)
}
