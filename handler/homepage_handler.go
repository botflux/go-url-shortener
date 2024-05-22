package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

type HomepageHandler struct {
	HomepageTemplate *template.Template
}

func (handler *HomepageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := handler.HomepageTemplate.Execute(w, make(map[string]interface{}))

	if err != nil {
		fmt.Fprint(w, "Something went wrong when rendering the template")
	}
}
