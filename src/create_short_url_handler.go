package main

import (
	"fmt"
	"html/template"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
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

	urlMapping := UrlMappingDocument{
		ShortenUrl:  RandomId(5),
		CompleteUrl: urlToShorten,
	}

	_, err := handler.Collection.InsertOne(r.Context(), urlMapping)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Cannot insert in the DB")
		return
	}

	w.WriteHeader(http.StatusCreated)
	handler.HomepageTemplate.Execute(w, urlMapping)
}
