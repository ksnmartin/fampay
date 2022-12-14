package cron

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ksnmartin/fampay/db"
	"go.mongodb.org/mongo-driver/mongo"
)

func Job(DB *mongo.Client) {
	response, err := http.Get("https://youtube.googleapis.com/youtube/v3/search?part=snippet&q=dogs&key=AIzaSyAcfQ9EMnnKBf0zj9eFzLRUKvHTJE09hWM")
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)
	items := data["items"].([]interface{})
	var snippetList []interface{}
	for _, obj := range items {
		snippet := obj.(map[string]interface{})["snippet"]
		snippetList = append(snippetList, snippet)
	}
	_, err = db.InsertData(DB, snippetList)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(time.Now(), ": data insert succesfull")
	}
}
