package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Сервер запущен")
}

func main() {
	http.HandleFunc("/notes", handler)
	fmt.Println("Сервер запущен на http://localhost:8080/notes")
	http.ListenAndServe(":8080", nil)
}
