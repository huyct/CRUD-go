package routes

import (
	"context"
	"fmt"
	"net/http"
	_ "reflect"

	"github.com/asaskevich/govalidator"
	jwt "github.com/huyct/CRUD-go/auth"
	db "github.com/huyct/CRUD-go/database"
	"github.com/huyct/CRUD-go/models"
	res "github.com/huyct/CRUD-go/utils"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	collection := db.ConnectPosts()

	var result []bson.M
	data, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		res.JSON(w, 500, "Internal Server Error")
		return
	}

	defer data.Close(context.Background())
	for data.Next(context.Background()) {
		var elem bson.M
		err := data.Decode(&elem)

		if err != nil {
			res.JSON(w, 500, "Internal Server failed")
			return
		}

		result = append(result, elem)
	}

	res.JSON(w, 200, result)
}

func GetMyPosts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username, err := jwt.EctractUsernameFromToken(r)

	if err != nil {
		res.JSON(w, 500, "Internal Server failed")
		return
	}

	collection := db.ConnectPosts()

	var result []bson.M
	data, err := collection.Find(context.Background(), bson.M{"creater": username})

	defer data.Close(context.Background())
	for data.Next(context.Background()) {
		var elem bson.M
		err := data.Decode(&elem)

		if err != nil {
			res.JSON(w, 500, "Internal Server failed")
			return
		}

		result = append(result, elem)
	}

	res.JSON(w, 200, result)
}

func CreatePost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	creater, err := jwt.EctractUsernameFromToken(r)

	if err != nil {
		res.JSON(w, 500, "Internal Server failed")
		return
	}

	title := r.PostFormValue("title")

	if govalidator.IsNull(title) {
		res.JSON(w, 400, "Data can not empty")
		return
	}

	title = models.Santize(title)
	uid := uuid.NewV4()

	id := fmt.Sprintf("%x-%x-%x-%x-%x", uid[0:4], uid[4:6], uid[6:8], uid[8:10], uid[10:])
	collection := db.ConnectPosts()

	newPost := bson.M{"id": id, "creater": creater, "title": title}

	_, errs := collection.InsertOne(context.TODO(), newPost)

	if errs != nil {
		res.JSON(w, 500, "Create post has failed")
		return
	}

	res.JSON(w, 200, "Create Successfully")
}

func EditPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	title := r.PostFormValue("title")
	username, err := jwt.EctractUsernameFromToken(r)

	if err != nil {
		res.JSON(w, 500, "Internal Server failed")
		return
	}

	if govalidator.IsNull(title) {
		res.JSON(w, 400, "Data can not empty")
		return
	}

	collection := db.ConnectPosts()

	var result bson.M
	errFind := collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&result)

	if errFind != nil {
		res.JSON(w, 404, "Post not found")
		return
	}

	creater := fmt.Sprintf("%v", result["creater"])

	if username != creater {
		res.JSON(w, 403, "Permission denied")
		return
	}

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"title": title}}

	_, errUpdate := collection.UpdateOne(context.TODO(), filter, update)

	if errUpdate != nil {
		res.JSON(w, 500, "Edit has failed")
		return
	}

	res.JSON(w, 200, "Edit Successfully")
}

func DeletePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	username, err := jwt.EctractUsernameFromToken(r)
	collection := db.ConnectPosts()

	if err != nil {
		res.JSON(w, 500, "Internal Server failed")
		return
	}

	var result bson.M
	errFind := collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&result)

	if errFind != nil {
		res.JSON(w, 404, "Post not found")
		return
	}

	creater := fmt.Sprintf("%v", result["creater"])

	if username != creater {
		res.JSON(w, 403, "Permission denied")
		return
	}

	errDelete := collection.FindOneAndDelete(context.TODO(), bson.M{"id": id}).Decode(&result)

	if errDelete != nil {
		res.JSON(w, 500, "Denied has failed")
		return
	}

	res.JSON(w, 200, "Delete post successfully")
}
