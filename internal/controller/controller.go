package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/travas-io/travas-op/internal/config"
	"github.com/travas-io/travas-op/internal/query"
	"go.mongodb.org/mongo-driver/mongo"
)

type Operator struct {
	App *config.Tools
	DB  query.Repo
}

func NewOperator(app *config.Tools, db *mongo.Client) *Operator {
	return &Operator{
		App: app,
		DB:  query.NewOperatorDB(app, db),
	}
}

func (o Operator) Home() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
