package query

import (
	"github.com/travas-io/travas-op/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repo interface {
	InsertUser(user *model.Operator) (int, primitive.ObjectID, error)
	VerifyUser(userID primitive.ObjectID) (bool, error)
	UpdateInfo(userID primitive.ObjectID, tk map[string]string) (bool, error)

	InsertPackage(tour *model.Tour) (bool, error)
	LoadTour(tourID primitive.ObjectID) (primitive.M, error)

	ValidTourRequest() ([]primitive.M, error)

	InsertTourGuide(tg *model.TourGuide) (bool, error)
	FindTourGuide(operatorID primitive.ObjectID) ([]primitive.M, error)
	UpdateTourGuide(guideID string) error
}
