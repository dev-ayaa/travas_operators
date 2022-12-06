package query

import (
	"context"
	"github.com/travas-io/travas-op/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// InsertPackage : this will help the operator to add new create tour packages in the
// tours collection with respect to a particular operators ID
func (op *OperatorDB) InsertPackage(tour *model.Tour) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	filter := bson.D{{Key: "destination", Value: tour.Destination}, {Key: "operator_id", Value: tour.OperatorID}}
	var res bson.M
	err := TourData(op.DB, "tours").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err := TourData(op.DB, "tours").InsertOne(ctx, tour)
			if err != nil {
				return false, err
			}
			return true, nil
		}
		op.App.ErrorLogger.Fatalf("error while searching for data : %v ", err)
	}

	return true, nil
}

func (op *OperatorDB) LoadTour(tourID primitive.ObjectID) (primitive.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	filter := bson.D{{Key: "_id", Value: tourID}}
	var res bson.M
	err := TourData(op.DB, "tours").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		op.App.ErrorLogger.Fatalf("error while searching for data : %v ", err)
	}
	return res, nil
}
