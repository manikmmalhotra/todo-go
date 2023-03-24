package router

import (
	"todo/controller"

	"github.com/gorilla/mux"
)

// func Router() http.Handler {
// 	router := chi.NewRouter()
// 	router.Use(middleware.Nosurf)
// 	router.Get("/", controller.HomeHandler)
// 	router.Mount("/todo", todoHandler())
// 	router.Delete("/todo/delete-completed", controller.DeleteCompleted)

// 	return router
// }

// func todoHandler() http.Handler {
// 	rg := chi.NewRouter()
// 	rg.Group(func(r chi.Router) {
// 		r.Get("/", controller.FetchTodods)
// 		r.Post("/", controller.CreateTodo)
// 		r.Put("/{id}", controller.UpdateTodo)
// 		r.Delete("/{id}", controller.DeleteOneTodo)
// 	})

// 	return rg
// }

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", controller.HomeHandler).Methods("GET")
	router.HandleFunc("/todo", controller.FetchTodods).Methods("GET")
	router.HandleFunc("/todo", controller.CreateTodo).Methods("POST")
	router.HandleFunc("/todo/{id}", controller.UpdateTodo).Methods("PUT")
	router.HandleFunc("/todo/{id}", controller.DeleteOneTodo).Methods("DELETE")
	router.HandleFunc("/todo/deleteall", controller.DeleteCompleted).Methods("DELETE")

	return router
}
