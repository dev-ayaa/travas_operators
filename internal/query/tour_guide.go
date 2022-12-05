package query

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (op *OperatorDB) FindTourGuide(operatorID primitive.ObjectID) ([]primitive.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	filter := bson.D{{Key: "operator_id", Value: operatorID}}

	cursor, err := TourData(op.DB, "tour_guide").Find(ctx, filter)
	if err != nil {
		op.App.ErrorLogger.Fatalf("error while searching for data : %v \n", err)

	}
	var res []bson.M
	if err := cursor.All(ctx, &res); err != nil {
		op.App.ErrorLogger.Fatal(err)
	}
	return res, nil
}
