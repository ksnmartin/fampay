package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/ksnmartin/fampay/cron"
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

	router.HandleFunc("/health", app.Health)
	router.HandleFunc("/search", app.SearchAPI)

	app.Server.Handler = router
}
func (app *App) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	jsonDocument := map[string]string{"status": "ok"}
	response, _ := json.Marshal(jsonDocument)
	w.Write(response)
}

func (app *App) SearchAPI(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	limitStr := r.URL.Query().Get("limit")
	limit := 25
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	} else if os.Getenv("DEFAULT_RESULTS") != "" {
		limit, _ = strconv.Atoi(os.Getenv("DEFAULT_RESULTS"))
	}
	row, err := db.SearchData(app.DB, q, limit)
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
	runCronJob := os.Getenv("CRON_INSTANCE")
	if runCronJob == "" || runCronJob == "true" {
		interval, _ := strconv.Atoi(os.Getenv("MINING_INTERVAL_IN_MIN"))
		app.Scheduler.Every(interval).Minute().Do(func() {
			cron.MiningCronJob(app.DB)
		})
		app.Scheduler.StartAsync()
	}

}

func CreateApp() *App {
	DB, err := db.Connect()
	if err != nil {
		log.Println("Database connection failed \n=>", err.Error())
	}
	srv := &http.Server{
		Addr: "0.0.0.0:" + os.Getenv("PORT"),
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
