package app

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"html/template"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"
)

// Router : Init router function
func Router(mux *mux.Router, database *mongo.Database) *mux.Router {
	api := API{Database: database}

	mux.HandleFunc("/api/create", api.createShort).Methods("POST")
	mux.HandleFunc("/", api.Index).Methods("GET")
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("public/static/"))))
	mux.HandleFunc("/{id}", api.parse).Methods("GET")

	return mux
}

type API struct {
	Database *mongo.Database
}

func (api API) parse(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	count, _ := api.Database.Collection("links").CountDocuments(context.TODO(), bson.M{"vanity": id})
	if count == 0 {
		t := template.Must(template.ParseFiles(
			"public/templates/404.html",
			"public/templates/partials/meta.html"))

		_ = t.ExecuteTemplate(w, "404", nil)

		return
	}

	var queryUser struct {
		Origin string `bson:"origin,omitempty"`
	}

	if err := api.Database.Collection("links").FindOne(context.TODO(), bson.M{"vanity": id}).Decode(&queryUser); err != nil {
		writeResp(w, 500, "Error!")
		return
	}

	http.Redirect(w, r, queryUser.Origin, http.StatusSeeOther)

	defer r.Body.Close()
}

func (api API) createShort(w http.ResponseWriter, r *http.Request) {
	var responseData struct {
		Link string `json:"link"`
	}

	err := json.NewDecoder(r.Body).Decode(&responseData)
	if err != nil {
		writeResp(w, 400, "Invalid JSON!")
		return
	}

	_, err = url.ParseRequestURI(responseData.Link)
	if err != nil {
		writeResp(w, 400, "Invalid URL!")
		return
	}

	random := randomString(5)

	count, _ := api.Database.Collection("links").CountDocuments(context.TODO(), bson.M{"vanity": random})
	for count != 0 {
		random = randomString(5)
		count, _ = api.Database.Collection("links").CountDocuments(context.TODO(), bson.M{"vanity": random})
	}

	_, _ = api.Database.Collection("links").InsertOne(context.TODO(), bson.M{"origin": responseData.Link, "vanity": random})

	type sendRes struct {
		URL string `json:"url"`
	}

	response, _ := json.Marshal(sendRes{URL: random}.URL)

	writeResp(w, 200, string(response))

	defer r.Body.Close()
}

// Stupid thingy for random strings
func randomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Write response function
func writeResp(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(msg))
}