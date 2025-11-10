package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"notes-api/internal/handlers"
	"notes-api/internal/storage"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Не удалось загрузить .env — используем системные переменные")
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
	r.HandleFunc("/notes/{id}", h.GetNoteByIDHandler).Methods("GET")
	r.HandleFunc("/notes/{id}", h.UpdateNoteHandler).Methods("PUT")
	r.HandleFunc("/notes/{id}", h.DeleteNoteHandler).Methods("DELETE")

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("Сервер запущен на http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
