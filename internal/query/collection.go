package query

import "go.mongodb.org/mongo-driver/mongo"

func OperatorData(db *mongo.Client, collection string) *mongo.Collection {
	var operators = db.Database("travas").Collection("operators")
	return operators
}

func TourData(db *mongo.Client, collection string) *mongo.Collection {
	var tours = db.Database("travas").Collection("tours")
	return tours
}

func TourGuideData(db *mongo.Client, collection string) *mongo.Collection {
	var guide = db.Database("travas").Collection("tour_guides")
	return guide
}
