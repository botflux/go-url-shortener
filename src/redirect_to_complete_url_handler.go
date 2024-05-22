package main

import (
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RedirectToCompleteURLHandler struct {
	Collection *mongo.Collection
}

func (handler *RedirectToCompleteURLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Unknown shorten URL")
		return
	}

	var mapping UrlMappingDocument

	err := handler.Collection.FindOne(r.Context(), bson.D{{Key: "shortenurl", Value: id}}).Decode(&mapping)

	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Shorten URL not found in DB")
		return
	}

	w.Header().Add("Location", mapping.CompleteUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
