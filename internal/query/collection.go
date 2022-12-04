package query

import "go.mongodb.org/mongo-driver/mongo"

func OperatorData(db *mongo.Client, collection string) *mongo.Collection {
	var operator = db.Database("travas").Collection("operator")
	return operator
}
