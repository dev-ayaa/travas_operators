package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/travas-io/travas-op/internal/controller"
)

func Routes(r *gin.Engine, o controller.Operator) {
	router := r.Use(gin.Logger(), gin.Recovery())
	router.Use(cors.Default())
	cookieData := cookie.NewStore([]byte("travas"))
	router.Use(sessions.Sessions("session", cookieData))
	router.GET("/", o.Welcome())
	//router.GET("/register", o.Register())
	router.POST("/register", o.ProcessRegister())
	//router.GET("/login", o.LoginPage())
	router.POST("/login", o.ProcessLogin())
	authRouter := r.Group("/auth")
	authRouter.Use(Authorization())
	{
		authRouter.POST("/verify/documents", o.VerifyDocument())

		authRouter.GET("/create/packages", o.TourPackagePage())
		authRouter.POST("/add/packages", o.ProcessTourPackage())

		authRouter.GET("/guide/load", o.GetTourGuide())
		authRouter.GET("/guide/select/assign", o.SelectTourGuide())
		authRouter.POST("/guide/select/delete/:id", o.DeleteTourGuide())
		authRouter.GET("/tour/preview", o.PreviewTour())
	}

}
