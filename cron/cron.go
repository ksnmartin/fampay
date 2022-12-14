package cron

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ksnmartin/fampay/db"
	"go.mongodb.org/mongo-driver/mongo"
)

func Job(DB *mongo.Client) {
	//get data from YT API and update it to data base
	response, err := http.Get("https://youtube.googleapis.com/youtube/v3/search?part=snippet&q=dogs&API_KEY=AIzaSyAcfQ9EMnnKBf0zj9eFzLRUKvHTJE09hWM")
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	var data []interface{}
	_ = json.Unmarshal(body, &data)
	db.InsertData(DB, data)
}
