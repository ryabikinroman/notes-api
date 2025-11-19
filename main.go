package main

import (
	"log"
	"net/http"
	"os"

	"notes-api/internal/handlers"
	"notes-api/internal/storage"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Println("Не удалось загрузить .env:", err)
	}

	db, err := storage.NewDB()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	h := &handlers.Handler{DB: db}

	r := mux.NewRouter()

	r.HandleFunc("/notes", h.GetAllNotesHandler).Methods("GET")
	r.HandleFunc("/notes", h.CreateNoteHandler).Methods("POST")
	r.HandleFunc("/auth/login", h.LoginHandler).Methods("POST")
	r.HandleFunc("/notes/{id}", h.GetNoteByIDHandler).Methods("GET")
	r.HandleFunc("/notes/{id}", h.UpdateNoteHandler).Methods("PUT")
	r.HandleFunc("/notes/{id}", h.DeleteNoteHandler).Methods("DELETE")

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Printf("Сервер запущен на http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
