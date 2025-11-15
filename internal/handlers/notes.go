package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"notes-api/internal/models"
	"notes-api/internal/storage"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Handler struct {
	DB *sql.DB
}

func (h *Handler) GetAllNotesHandler(w http.ResponseWriter, r *http.Request) {

	notes, err := storage.GetAllNotes(h.DB)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (h *Handler) CreateNoteHandler(w http.ResponseWriter, r *http.Request) {

	var payload models.Note

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	id, err := storage.CreateNote(h.DB, &payload)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Заметка успешно создана",
		"id":      id,
	})
}

func (h *Handler) GetNoteByIDHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID не указан", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	note, err := storage.GetNoteByID(h.DB, id)
	if err != nil {
		if strings.Contains(err.Error(), "не найдена") {
			http.Error(w, "Запись не найдена", http.StatusNotFound)
		} else {
			WriteError(w, http.StatusInternalServerError, "Ошибка сервера")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (h *Handler) UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID не указан", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	var payload models.Note
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	err = storage.UpdateNote(h.DB, id, &payload)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Заметка не найдена", http.StatusNotFound)
			return
		}
		WriteError(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Заметка с ID=%d успешно обновлена", id),
	})
}

func (h *Handler) DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID не указан", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	err = storage.DeleteNote(h.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Заметка не найдена", http.StatusNotFound)
			return
		}
		WriteError(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}
}
