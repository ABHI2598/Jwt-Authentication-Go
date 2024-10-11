package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/abhinavpandey/jwtProject/database"
	"github.com/abhinavpandey/jwtProject/helpers"
	"github.com/abhinavpandey/jwtProject/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection = database.OpenCollection(database.Client, "users")
var validate = validator.New()

func HashPassword(password string) string {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal(err)
	}

	return string(pass)
}

func VerifyPassword(incomingPassword string, dbPassWord string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(dbPassWord), []byte(incomingPassword))
	isValid := true
	msg := ""

	if err != nil {
		log.Fatal("Password Mismatch")
		msg = "Password is incorrect"
		isValid = false
	}

	return isValid, msg
}

func SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		var user models.User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count1, err := userCollection.CountDocuments(c, bson.M{"email": user.Email})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
		}

		if count1 > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Email or phone already exists"})
			return
		}

		count2, err := userCollection.CountDocuments(c, bson.M{"phone": user.Phone})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
		}

		if count2 > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Email or phone already exists"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()

		token, refreshToken, _ := helpers.GenerateToken(*user.Email, *user.FirstName, *user.LastName, *user.UserType, *&user.UserId)

		user.Token = &token
		user.RefreshToken = &refreshToken

		ctx.SetCookie("token", *user.Token, time.Now().Hour(), "", "localhost", false, true)

		result, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "InsertiON fAILED"})
		}

		ctx.JSON(http.StatusCreated, result)

	}
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Wrong input"})
			return
		}

		filter := bson.M{"email": user.Email}
		err := userCollection.FindOne(c, filter).Decode(&foundUser)
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "No docuement found"})
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		isValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !isValid {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}

		if foundUser.Email == nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		}

		token, refreshToken, _ := helpers.GenerateToken(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, *foundUser.UserType, *&foundUser.UserId)

		helpers.UpdateAllTokens(token, refreshToken, foundUser.UserId)

		objectId, _ := primitive.ObjectIDFromHex(foundUser.UserId)
		err = userCollection.FindOne(c, bson.M{"_id": objectId}).Decode(&foundUser)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.SetCookie("token", *foundUser.Token, time.Now().Local().Hour(), "", "localhost", false, true)
		ctx.JSON(http.StatusOK, foundUser)

	}
}

func GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := helpers.CheckUserType(ctx, "ADMIN"); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		}

		var c, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(ctx.Query("recordPerPage"))
		if err != nil && recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(ctx.Query("page"))
		if err1 != nil && page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage

		pipeline := []bson.M{
			{"$match": bson.M{}},
			{"$skip": startIndex},
			{"$limit": recordPerPage},
			{"$group": bson.M{"_id": nil, "data": bson.M{"$push": "$$ROOT"}, "total_users": bson.M{"$sum": 1}}},
			{"$project": bson.M{"_id": 0, "user_items": bson.M{"$slice": []interface{}{"$data", startIndex, recordPerPage}}, "total_users": 1}},
		}

		result, err := userCollection.Aggregate(c, pipeline)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while listing items"})
		}

		if result == nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Aggregation result is nil"})
			return
		}

		var allUsers []bson.M
		if err = result.All(c, &allUsers); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured pushing items"})
		}
		ctx.JSON(http.StatusOK, allUsers[0])

	}
}

func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("id")

		if err := helpers.MatchUserTypeToUid(ctx, userId); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		}

		c, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()
		var user models.User
		objectId, _ := primitive.ObjectIDFromHex(userId)

		filter := bson.M{"_id": objectId}
		err := userCollection.FindOne(c, filter).Decode(&user)
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "No docuement found"})
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, user)
	}
}
