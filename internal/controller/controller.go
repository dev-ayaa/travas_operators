package controller

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/travas-io/travas-op/internal/config"
	"github.com/travas-io/travas-op/internal/encrypt"
	"github.com/travas-io/travas-op/internal/query"
	"github.com/travas-io/travas-op/internal/token"
	"github.com/travas-io/travas-op/model"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strings"
	"time"
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

// Welcome : This method render the welcome page of the flutter mobile application
func (op *Operator) Welcome() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Todo : render the home page of the application
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

// Register : this Handler will render and show the register page for user
func (op *Operator) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

// ProcessRegister : As the name implies , this method will help to process all the registration process
// of the user
func (op *Operator) ProcessRegister() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user model.Operator

		if err := ctx.Request.ParseForm(); err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
		}
		user.Email = ctx.Request.Form.Get("email")
		user.Phone = ctx.Request.Form.Get("phone")
		user.Password = ctx.Request.Form.Get("password")
		user.CompanyName = ctx.Request.Form.Get("company_name")
		user.ConfirmPassword = ctx.Request.Form.Get("check_password")
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Token = ""
		user.NewToken = ""
		user.ToursList = []model.Tour{}

		if user.Password != user.ConfirmPassword {
			_ = ctx.AbortWithError(http.StatusInternalServerError, errors.New("passwords did not match"))
		}

		user.Password, _ = encrypt.Hash(user.Password)
		user.ConfirmPassword, _ = encrypt.Hash(user.ConfirmPassword)

		if err := op.App.Validator.Struct(&user); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); !ok {
				_ = ctx.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
				log.Println(err)
				return
			}
		}

		track, userID, err := op.DB.InsertUser(user)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("error while adding new user"))
			return
		}
		cookieData := sessions.Default(ctx)

		userInfo := model.UserInfo{
			ID:       userID,
			Email:    user.Email,
			Password: user.Password,
		}
		cookieData.Set("info", userInfo)

		if err := cookieData.Save(); err != nil {
			log.Println("error from the session storage")
			_ = ctx.AbortWithError(http.StatusNotFound, gin.Error{Err: err})
			return
		}
		switch track {
		case 1:
			// add the user id to session
			// redirect to the home page of the application
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Existing Account, Go to the Login page",
			})
		case 0:
			//	after inserting new user to the database
			//  notify the user to verify their  details via mail
			//  OR
			//  Send notification message on the page for them to login
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Registered Successfully",
			})
		}
	}
}

// LoginPage : this will show the login page for user
func (op *Operator) LoginPage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSONP(http.StatusOK, gin.H{})
	}
}

// ProcessLogin : this method will help to parse, verify, and as well as authenticate the user
// login details, and it also helps to generate jwt token for restricted routers

func (op *Operator) ProcessLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := ctx.Request.ParseForm(); err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
		}
		email := ctx.Request.Form.Get("email")
		password := ctx.Request.Form.Get("password")

		cookieData := sessions.Default(ctx)
		userInfo := cookieData.Get("info").(model.UserInfo)

		verified, err := encrypt.Verify(password, userInfo.Password)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, errors.New("cannot verify user input password"))
		}
		if verified {
			switch {
			case email == userInfo.Email:
				_, checkErr := op.DB.VerifyUser(userInfo.ID)

				if checkErr != nil {
					_ = ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("unregistered user %v", checkErr))
				}
				// generate the jwt token
				t1, t2, err := token.Generate(userInfo.Email, userInfo.ID)
				if err != nil {
					_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token no generated : %v ", err))
				}

				var tk map[string]string
				tk = map[string]string{"t1": t1, "t2": t2}

				// update the database adding the token to user database
				_, updateErr := op.DB.UpdateInfo(userInfo.ID, tk)
				if updateErr != nil {
					_ = ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("unregistered user %v", updateErr))
				}

				ctx.SetCookie("authorization", t1, 60*60*24*7, "/", "localhost", false, true)
				ctx.JSONP(http.StatusOK, gin.H{"message": "Welcome to user homepage"})
			}
		}
	}
}

func (op *Operator) VerifyDocument() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//	todo --> verify document upload by the tour operator
		//	this will involve scanning the pdf format of the document
		//	signature and other details needed
	}
}

func (op *Operator) TourPackagePage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSONP(http.StatusOK, gin.H{})
	}
}

func (op *Operator) ProcessTourPackage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var tour model.Tour
		if err := ctx.Request.ParseForm(); err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
		}

		cookieData := sessions.Default(ctx)
		userInfo := cookieData.Get("info").(model.UserInfo)

		tour.OperatorID = userInfo.ID
		tour.Title = ctx.Request.Form.Get("title")
		tour.Destination = strings.TrimSpace(strings.ToLower(ctx.Request.Form.Get("destination")))
		tour.MeetingPoint = ctx.Request.Form.Get("meeting_point")
		tour.StartTime = ctx.Request.Form.Get("start_time")
		tour.StartDate = ctx.Request.Form.Get("start_date")
		tour.Price = ctx.Request.Form.Get("price")
		tour.Language = ctx.Request.Form.Get("language")
		tour.Description = ctx.Request.Form.Get("description")
		tour.NumberOfTourist = ctx.Request.Form.Get("number_of_tourists")
		tour.WhatToExpect = append(tour.WhatToExpect, ctx.Request.Form.Get("what_to_expect"))
		tour.Rules = append(tour.Rules, ctx.Request.Form.Get("rules"))

		if err := op.App.Validator.Struct(&tour); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); !ok {
				_ = ctx.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
				log.Println(err)
				return
			}
		}
		_, err := op.DB.InsertPackage(tour)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
			return
		}

		ctx.JSONP(http.StatusCreated, gin.H{"Message": "new tour package created"})

	}
}

func (op *Operator) GetTourGuide() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookieData := sessions.Default(ctx)
		userInfo := cookieData.Get("info").(model.UserInfo)

		arrRes, err := op.DB.FindTourGuide(userInfo.ID)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
		}

		ctx.JSONP(http.StatusOK, gin.H{"TourGuides": arrRes})
	}
}
