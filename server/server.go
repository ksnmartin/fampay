package server

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ksnmartin/fampay/db"
)

type App struct {
	DB     *sql.DB
	Server *http.Server
}

func (app *App) AddRoutes() {
	router := mux.NewRouter()

	router.HandleFunc("/", app.GetYoutubeData)
	router.HandleFunc("/search", app.SearchAPI)

	app.Server.Handler = router
}
func (app *App) GetYoutubeData(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("limit")
	limit := 25
	if q != "" {
		limit, _ = strconv.Atoi(q)
	}
	_, _ := db.GetAllData(app.DB, limit)
}
func (app *App) SearchAPI(http.ResponseWriter, *http.Request) {

}

func (app *App) AddCronJob(http.ResponseWriter, *http.Request) {

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
		DB:     DB,
		Server: srv,
	}
	app.AddRoutes()
	return &app
}
