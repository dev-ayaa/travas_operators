package main

import (
	"github.com/travas-io/travas-op/internal/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, err := ctx.Request.Cookie("authorization")
		if err != nil {
			if err == http.ErrNoCookie {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}
		tokenString := c.Value

		if tokenString == "" {
			_ = ctx.AbortWithError(http.StatusNoContent, errors.New("no value for authorization header"))
		}

		parse, err := token.Parse(tokenString)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, gin.Error{Err: err})
		}
		ctx.Set("pass", tokenString)
		ctx.Set("id", parse.ID)
		ctx.Set("email", parse.Email)
		ctx.Next()
	}
}
