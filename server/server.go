package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/ksnmartin/fampay/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	DB        *mongo.Client
	Server    *http.Server
	Scheduler *gocron.Scheduler
}

func (app *App) AddRoutes() {
	router := mux.NewRouter()

	router.HandleFunc("/data", app.GetYoutubeData)
	router.HandleFunc("/search", app.SearchAPI)

	app.Server.Handler = router
}
func (app *App) GetYoutubeData(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("limit")
	limit := 25
	if q != "" {
		limit, _ = strconv.Atoi(q)
	}
	row, err := db.GetData(app.DB, int64(limit))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var jsonDocuments []map[string]interface{}
	var bsonDocument bson.D
	var jsonDocument map[string]interface{}
	var temporaryBytes []byte
	for row.Next(context.Background()) {
		row.Decode(&bsonDocument)
		temporaryBytes, _ = bson.MarshalExtJSON(bsonDocument, true, true)
		_ = json.Unmarshal(temporaryBytes, &jsonDocument)
		jsonDocuments = append(jsonDocuments, jsonDocument)
	}
	response, _ := json.Marshal(jsonDocuments)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (app *App) SearchAPI(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	row, err := db.SearchData(app.DB, q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
	var jsonDocuments []map[string]interface{}

	for row.Next(context.TODO()) {
		var bsonDocument bson.D
		var jsonDocument map[string]interface{}
		var temporaryBytes []byte
		row.Decode(&bsonDocument)
		temporaryBytes, _ = bson.MarshalExtJSON(bsonDocument, true, true)
		_ = json.Unmarshal(temporaryBytes, &jsonDocument)
		fmt.Println(jsonDocument)
		jsonDocuments = append(jsonDocuments, jsonDocument)
	}
	fmt.Println(jsonDocuments)
	response, _ := json.Marshal(jsonDocuments)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (app *App) AddCronJob() {
	fmt.Println("set up")
	app.Scheduler.Every(30).Second().Do(func() {
		//cron.Job(app.DB)
		fmt.Println("called at : ", time.Now())
	})
	app.Scheduler.StartAsync()
	app.Scheduler.RunAll()

}

func CreateApp() *App {
	DB, err := db.Connect()
	if err != nil {
		log.Println("Database connection failed \n=>", err.Error())
	}
	srv := &http.Server{
		Addr: "localhost:8000",
	}
	app := App{
		DB:        DB,
		Server:    srv,
		Scheduler: gocron.NewScheduler(time.UTC),
	}
	app.AddRoutes()
	app.AddCronJob()
	return &app
}
