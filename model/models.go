package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Data struct {
	ResourceId string               `bson:"_id,omitempty"`
	TotalTime  primitive.Decimal128 `bson:"totalTime"`
}

type NewData struct {
	ResourceId string `bson:"resourceID"`
	TimeSpent  string `bson:"timeSpent"`
	StartTime  string `bson:"startDateTime"`
	EndTime    string `bson:"endDateTime"`
}
