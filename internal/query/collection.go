package query

import "go.mongodb.org/mongo-driver/mongo"

func OperatorData(db *mongo.Client, collection string) *mongo.Collection {
	var operators = db.Database("travasdb").Collection(collection)
	return operators
}

func TourData(db *mongo.Client, collection string) *mongo.Collection {
	var tours = db.Database("travasdb").Collection(collection)
	return tours
}

func TourGuideData(db *mongo.Client, collection string) *mongo.Collection {
	var guide = db.Database("travasdb").Collection(collection)
	return guide
}

func TouristsData(db *mongo.Client, collection string) *mongo.Collection {
	var touristData = db.Database("travasdb").Collection(collection)
	return touristData
}
