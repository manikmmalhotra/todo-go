package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"todo/helper"
	"todo/model"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/justinas/nosurf"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var rnd *renderer.Render
var db *mongo.Collection
var ctx = context.TODO()

func init() {
	rnd = renderer.New()
	if err := godotenv.Load(); err != nil {
		log.Fatal("There is no env file")
	}
	mongoUri := os.Getenv("MONGO_URI")
	fmt.Println(mongoUri)
	clientOptions := options.Client().ApplyURI("mongodb+srv://minks:7418529633333@netflix.fnsy5mv.mongodb.net/?retryWrites=true&w=majority")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	db = client.Database("todo").Collection("task")

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("this is working ")
	todos, err := helper.FetchTododsFormDB(db)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todo",
			"error":   err,
		})
		return
	}

	data := model.TemplateData{
		CSRFToken: nosurf.Token(r),
		Todos:     todos,
	}

	err = rnd.Template(w, http.StatusOK, []string{""}, data)
	if err != nil {
		log.Fatal(err)
	}

}

func FetchTodods(w http.ResponseWriter, r *http.Request) {
	todos, err := helper.FetchTododsFormDB(db)

	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todo",
			"error":   err,
		})

		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"data": todos,
	})
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {

	var t model.TodoModel
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusProcessing, err)
		return
	}

	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title is required",
		})

		return
	}

	tm := model.TodoModel{
		ID:        t.ID,
		Title:     t.Title,
		Completed: t.Completed,
		CreatedAt: time.Now(),
	}

	_, err := db.InsertOne(ctx, &tm)

	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to save",
			"error":   err,
		})
		return
	}

	todos, err := helper.FetchTododsFormDB(db)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todo",
			"error":   err,
		})

		return
	}

	rnd.JSON(w, http.StatusCreated, renderer.M{
		"message": "todo is created successfully",
		"todos":   todos,
	})
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	// id := strings.TrimSpace(chi.URLParam(r, "id"))

	params := mux.Vars(r)

	objID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The id is invalid",
		})
		return
	}

	var t model.TodoModel

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusProcessing, err)
		return
	}

	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title field is requried",
		})
		return
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"title": t.Title, "completed": t.Completed}}

	_, err = db.UpdateOne(ctx, filter, update)

	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to update todo",
			"error":   err,
		})
		return
	}

	todos, err := helper.FetchTododsFormDB(db)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todo",
			"error":   err,
		})

		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Todo updated successfully",
		"todos":   todos,
	})

}

func DeleteOneTodo(w http.ResponseWriter, r *http.Request) {
	// id := strings.TrimSpace(chi.URLParam(r, "id"))
	params := mux.Vars(r)

	fmt.Println(params["id"])

	objID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The id is invalid",
		})
		return
	}

	_, err = db.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to delete todo",
			"error":   err,
		})
		return
	}

	todos, err1 := helper.FetchTododsFormDB(db)
	if err1 != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todo",
			"error":   err1,
		})

		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "todo is successfully deleted",
		"todos":   todos,
	})
}

func DeleteCompleted(w http.ResponseWriter, r *http.Request) {
	filter := bson.M{
		"completed": bson.M{
			"$eq": true,
		},
	}
	_, err := db.DeleteMany(ctx, filter)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to delete completed todos",
			"error":   err,
		})
		return
	}

	todos, err1 := helper.FetchTododsFormDB(db)
	if err1 != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todo",
			"error":   err1,
		})

		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "todo is successfully deleted",
		"todos":   todos,
	})
}
