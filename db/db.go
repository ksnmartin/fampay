package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Client, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://martin:mishravikas@cluster0.k1p632w.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client, err
}
func GetData(dataBase *mongo.Client, limit int64) (*mongo.Cursor, error) {
	collection := dataBase.Database("Youtube").Collection("searchResult")
	filter := bson.D{{}}
	options := options.Find()

	// Sort by `_id` field descending
	options.SetSort(bson.D{{Key: "publishTime", Value: -1}})

	// Limit by 10 documents only
	options.SetLimit(limit)
	result, err := collection.Find(context.TODO(), filter, options)
	return result, err
}

func InsertData(dataBase *mongo.Client, data []interface{}) (*mongo.InsertManyResult, error) {
	collection := dataBase.Database("Youtube").Collection("searchResult")
	result, err := collection.InsertMany(context.TODO(), data, options.InsertMany().SetOrdered(false))
	return result, err
}

func SearchData(dataBase *mongo.Client, query string) (*mongo.Cursor, error) {
	collection := dataBase.Database("Youtube").Collection("searchResult")
	filter := bson.A{
		bson.D{
			{"$search",
				bson.D{
					{"index", "default"},
					{"compound",
						bson.D{
							{"should",
								bson.A{
									bson.D{
										{"text",
											bson.D{
												{"query", query},
												{"path", "description"},
											},
										},
									},
									bson.D{
										{"text",
											bson.D{
												{"query", query},
												{"path", "title"},
											},
										},
									},
								},
							},
							{"minimumShouldMatch", 1},
						},
					},
				},
			},
		},
	}
	result, err := collection.Aggregate(context.TODO(), filter)
	return result, err
}
