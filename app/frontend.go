package app

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"html/template"
	"net/http"
)

// Index : This is the index page
func (api API) Index(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles(
		"public/templates/home.html",
		"public/templates/partials/meta.html"))

	count, _ := api.Database.Collection("links").CountDocuments(context.TODO(), bson.M{})

	_ = t.ExecuteTemplate(w, "home", count)
	defer r.Body.Close()
}