package query

import (
	"github.com/travas-io/travas-op/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
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

func (op OperatorDB) InsertOperator() {
	return
}
