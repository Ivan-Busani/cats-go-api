package main

import (
	"fmt"
	"log"
	"net/http"

	"cats-go-api/internal/database"
	"cats-go-api/internal/handler"
	"cats-go-api/internal/repository"
	"cats-go-api/internal/service"
)

func main() {
	db := database.MustNewPostgres()
	defer database.CloseDB(db)

	catRepo := repository.NewPostgresCatRepository(db)
	catSvc := service.NewCatService(catRepo)
	catHandler := handler.NewCatHandler(catSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.HealthHandler)
	mux.HandleFunc("GET /api/v1/cats/list", catHandler.List)
	mux.HandleFunc("GET /api/v1/cats/{id}", catHandler.GetByID)
	mux.HandleFunc("GET /api/v1/cats/cat_id/{cat_id}", catHandler.GetByCatID)
	mux.HandleFunc("POST /api/v1/cats/save", catHandler.Save)
	mux.HandleFunc("PUT /api/v1/cats/update/{id}", catHandler.Update)
	mux.HandleFunc("DELETE /api/v1/cats/delete/{id}", catHandler.Delete)

	fmt.Println("Server is running on port 8001")
	if err := http.ListenAndServe(":8001", mux); err != nil {
		log.Fatal(err)
	}
}
