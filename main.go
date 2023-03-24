package main

import (
	"fmt"
	"log"
	"net/http"
	"todo/router"
)

// func main() {
// 	srv := &http.Server{
// 		Addr:    ":9000",
// 		Handler: router.Router(),
// 	}
// 	err := srv.ListenAndServe()
// 	log.Println("Server is running on http://localhost:9000")
// 	if err != nil {
// 		panic(err)
// 	}
// }

func main() {
	fmt.Println("MongoDB API")
	r := router.Router()
	fmt.Println("Server is getting started...")
	log.Fatal(http.ListenAndServe(":4000", r))
	fmt.Println("Listening at port 4000 ...")
}
