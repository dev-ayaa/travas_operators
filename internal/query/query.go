package query

import (
	"context"
	"github.com/travas-io/travas-op/internal/config"
	"github.com/travas-io/travas-op/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type OperatorDB struct {
	App *config.Tools
	DB  *mongo.Client
}

func NewOperatorDB(app *config.Tools, db *mongo.Client) Repo {
	return &OperatorDB{
		App: app,
		DB:  db,
	}
}

// InsertUser : this will check for existing user and also insert new to the operator collection
func (op *OperatorDB) InsertUser(user model.Operator) (int, primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{{Key: "email", Value: user.Email}}

	var res bson.M
	err := OperatorData(op.DB, "operators").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			user.ID = primitive.NewObjectID()
			_, insertErr := OperatorData(op.DB, "operators").InsertOne(ctx, user)
			if insertErr != nil {
				op.App.ErrorLogger.Fatalf("cannot add user to the database : %v ", insertErr)
			}
			return 0, user.ID, nil
		}
		op.App.ErrorLogger.Fatal(err)
	}

	id := (res["_id"]).(primitive.ObjectID)

	return 1, id, err
}

// VerifyUser : this method will verify the user login details store in the database
// and compare with the input details
func (op *OperatorDB) VerifyUser(userID primitive.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var result bson.M

	filter := bson.D{{Key: "_id", Value: userID}}
	err := OperatorData(op.DB, "operators").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			op.App.ErrorLogger.Println("no document found for this query")
			return false, err
		}
		op.App.ErrorLogger.Fatalf("cannot execute the database query perfectly : %v ", err)
	}

	return true, nil
}

func (op *OperatorDB) UpdateInfo(userID primitive.ObjectID, tk map[string]string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: userID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "token", Value: tk["t1"]}, {Key: "new_token", Value: tk["t2"]}}}}

	_, err := OperatorData(op.DB, "operators").UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}
