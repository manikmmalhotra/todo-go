package helper

import (
	"context"
	"errors"
	"todo/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx = context.TODO()

func FetchTododsFormDB(db *mongo.Collection) ([]model.Todo, error) {
	var todos []model.TodoModel
	todolist := []model.Todo{}
	cur, err := db.Find(ctx, bson.D{})
	if err != nil {
		defer cur.Close(ctx)
		return todolist, errors.New("Failer to Fetch todo")

	}

	if err = cur.All(ctx, &todos); err != nil {
		return todolist, errors.New("Failed to Load Data")
	}

	for _, t := range todos {
		todolist = append(todolist, model.Todo{
			ID:        t.ID.Hex(),
			Title:     t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt,
		})
	}
	return todolist, nil
}
