// main.go
package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"myapp/models"
	"myapp/handlers"
	"myapp/middleware"
	"log"

	"github.com/rs/cors"
)

func main() {
	db := models.InitDB()
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/register", handlers.RegisterHandler(db)).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler(db)).Methods("POST")
	r.Handle("/change-password", middleware.AuthMiddleware(http.HandlerFunc(handlers.ChangePasswordHandler(db)))).Methods("POST")

	corsHandler := cors.Default().Handler(r)

	log.Println("Starting the server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
