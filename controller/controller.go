package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThePratikSah/gomongoapi/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var collection *mongo.Collection

func init() {
	URI := os.Getenv("URI")
	DB := os.Getenv("DB")
	COLLECTION := os.Getenv("COLLECTION")

	clientOptions := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB connection success")
	collection = client.Database(DB).Collection(COLLECTION)

	fmt.Println("Collection ref is ready")
}

func getAggCalculation() []primitive.M {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Aggregate(ctx, bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "startDateTime", Value: bson.D{{Key: "$gte", Value: "1567276200000"}}},
					{Key: "endDateTime", Value: bson.D{{Key: "$lte", Value: "1569868199000"}}},
				},
			},
		},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: "$resourceID"},
					{Key: "totalTime", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$toDecimal", Value: "$timeSpent"}}}}},
				},
			},
		},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
	})

	if err != nil {
		log.Fatal(err)
	}

	var data []bson.M

	for cursor.Next(ctx) {
		var singleData bson.M
		if err := cursor.Decode(&singleData); err != nil {
			log.Fatal(err)
		}

		data = append(data, singleData)
	}

	return data
}

func getWithoutAggr() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetProjection(bson.D{
		{Key: "_id", Value: 0},
		{Key: "resourceID", Value: 1},
		{Key: "timeSpent", Value: 1},
		{Key: "startDateTime", Value: 1},
		{Key: "endDateTime", Value: 1},
	})

	results, err := collection.Find(ctx, bson.D{
		{Key: "startDateTime", Value: bson.D{{Key: "$gte", Value: "1567276200000"}}},
		{Key: "endDateTime", Value: bson.D{{Key: "$lte", Value: "1569868199000"}}},
	}, opts)

	if err != nil {
		log.Fatal(err)
	}

	dataWithTime := make(map[string]string)
	for results.Next(ctx) {
		var singleObject model.NewData
		if err = results.Decode(&singleObject); err != nil {
			log.Fatal(err)
		}

		resourceId := singleObject.ResourceId
		timeSpent := singleObject.TimeSpent
		val, present := dataWithTime[resourceId]

		if !present {
			dataWithTime[resourceId] = timeSpent
		} else {
			newTime, _ := strconv.Atoi(timeSpent)
			existingTimeSpent, _ := strconv.Atoi(val)
			dataWithTime[resourceId] = strconv.Itoa(newTime + existingTimeSpent)
		}
	}

	return dataWithTime
}

func GetAggregationData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := getAggCalculation()
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatal(err)
	}
}

func GetDataWithAgg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := getWithoutAggr()
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatal(err)
	}
}
