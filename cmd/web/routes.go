package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/travas-io/travas-op/internal/controller"
)

func Routes(r *gin.Engine, t controller.Operator) {
	router := r.Use(gin.Logger(), gin.Recovery())
	router.Use(cors.Default())
	cookieData := cookie.NewStore([]byte("travas"))
	router.Use(sessions.Sessions("session", cookieData))

}

