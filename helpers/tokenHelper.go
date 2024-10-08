package helpers

import (
	"context"
	"os"
	"time"

	"github.com/abhinavpandey/jwtProject/database"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UserType  string
	Uid       string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")
var userCollection = database.OpenCollection(database.Client, "users")

func GenerateToken(email string, firstName string, lastName string, userType string, uid string) (token string, refreshToken string, err error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UserType:  userType,
		Uid:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err1 := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err1 != nil {
		panic(err1)
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		panic(err)
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		})

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "Invalid token"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Token is Expired"
	}

	return claims, msg
}

func UpdateAllTokens(token string, refreshToken string, userId string) {
	var c, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var updateObj primitive.D

	updateObj = append(updateObj, primitive.E{Key: "token", Value: token})
	updateObj = append(updateObj, primitive.E{Key: "refreshToken", Value: refreshToken})

	UpdatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, primitive.E{Key: "updated_at", Value: UpdatedAt})

	filter := bson.M{"user_id": userId}
	// update := bson.M {
	// 	"$set": bson.M {
	// 		"token": token,
	// 		"refreshToken": refreshToken,
	// 		"updated_at": UpdatedAt,
	// 	},
	// }

	_, err := userCollection.UpdateMany(c, filter, bson.M{"$set": updateObj})
	if err != nil {
		panic(err)
	}
}
