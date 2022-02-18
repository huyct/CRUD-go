package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/huyct/CRUD-go/auth"
	"github.com/huyct/CRUD-go/database"
	"github.com/huyct/CRUD-go/models"
	"github.com/huyct/CRUD-go/utils"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
)

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	if govalidator.IsNull(username) || govalidator.IsNull(password) {
		res.JSON(w, 400, "Data can not empty")
	}

	username = models.Santize(username)
	password = models.Santize(password)

	collection := db.ConnectUsers()

	var result bson.M
	err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)

	if err != nil {
		res.JSON(w, 400, "Username or Password incorrect")
		return
	}

	//convert interface to string
	hashedPassword := fmt.Sprint("%v", result["password"])

	err = models.CheckPasswordHash(hashedPassword, password)

	if err != nil {
		res.JSON(w, 400, "Username or Password incorrect")
		return
	}

	token, errCreate := jwt.Create(username)

	if errCreate != nil {
		res.JSON(w, 500, "Internal Server Error")
		return
	}

	res.JSON(w, 200, token)
}

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	email := r.PostFormValue("email")

	if govalidator.IsNull(username) || govalidator.IsNull(password) || govalidator.IsNull(email) {
		res.JSON(w, 400, "Data cannot be empty")
		return
	}

	if !govalidator.IsEmail(email) {
		res.JSON(w, 400, "Email is invalid")
		return
	}

	username = models.Santize(username)
	email = models.Santize(email)
	password = models.Santize(password)

	collection := db.ConnectPosts()

	var result bson.M
	errFindUsername := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)
	errFindEmail := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&result)

	if errFindUsername == nil && errFindEmail == nil {
		res.JSON(w, 400, "User does exists")
		return
	}

	password, err := models.Hash(password)

	if err != nil {
		res.JSON(w, 500, "Regiter has failed")
		return
	}

	newUser := bson.M{"username": username, "email": email, "password": password}

	_, errs := collection.InsertOne(context.TODO(), newUser)

	if errs != nil {
		res.JSON(w, 500, "Register has failed")
		return
	}

	res.JSON(w, 200, "Register Successfully")
}
