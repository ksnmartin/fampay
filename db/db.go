package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Client, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(os.Getenv("DB_ADDRESS")).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client, err
}

func InsertData(dataBase *mongo.Client, data []interface{}) (*mongo.InsertManyResult, error) {
	collection := dataBase.Database("Youtube").Collection("searchResult")
	result, err := collection.InsertMany(context.TODO(), data, options.InsertMany().SetOrdered(false))
	//insert many with order set as false will insert unique values but return an error for no unique documents
	return result, err
}

func SearchData(dataBase *mongo.Client, query string, limit int) (*mongo.Cursor, error) {
	collection := dataBase.Database("Youtube").Collection("searchResult")
	paginationFilter := bson.D{{"$limit", limit}}
	sortFilter := bson.D{{"$sort", bson.D{{"publishedAt", -1}}}} //sorted in descending order
	filter := bson.A{}                                           //query pipline array
	if query != "" {
		partialTextSearchFilter := bson.D{
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
												{"path", "title"},
											},
										},
									},
									bson.D{
										{"text",
											bson.D{
												{"query", query},
												{"path", "description"},
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
		} //partial search text query
		filter = append(filter, partialTextSearchFilter)
	} else {
		//dont sort if query is present as search ranking will lose relavance
		filter = append(filter, sortFilter)
	}
	filter = append(filter, paginationFilter)
	result, err := collection.Aggregate(context.TODO(), filter)
	return result, err
}
