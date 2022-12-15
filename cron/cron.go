package cron

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ksnmartin/fampay/db"
	"go.mongodb.org/mongo-driver/mongo"
)

func MiningCronJob(DB *mongo.Client) {
	apiKeys := strings.Split(os.Getenv("API_KEY"), ",")
	publishedAfter := os.Getenv("PUBLISHED_AFTER")
	topic := os.Getenv("TOPIC")
	for _, apiKey := range apiKeys {
		//loop to use the next API key if previous key fails
		response, err := http.Get("https://youtube.googleapis.com/youtube/v3/search?part=snippet&q=" + topic + "&type=video&order=date&publishedAfter=" + publishedAfter + "&key=" + apiKey)
		if err != nil || response.StatusCode >= http.StatusBadRequest {
			continue
		} else if response.StatusCode == http.StatusOK {
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)
			var jsonDocument map[string]interface{}
			_ = json.Unmarshal(body, &jsonDocument)
			items := jsonDocument["items"].([]interface{})
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
			return
		}
	}

}
