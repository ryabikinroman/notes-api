package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"notes-api/internal/models"
	"notes-api/internal/storage"
	"strconv"

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
		WriteError(w, http.StatusBadRequest, "Неверный формат JSON")
		return
	}

	if err := payload.Validate(); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
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
	id, err := parseID(w, r)
	if err != nil {
		return
	}

	note, err := storage.GetNoteByID(h.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			WriteError(w, http.StatusNotFound, "Заметка не найдена")
		} else {
			WriteError(w, http.StatusInternalServerError, "Ошибка сервера")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (h *Handler) UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(w, r)
	if err != nil {
		return
	}

	var payload models.Note
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteError(w, http.StatusBadRequest, "Неверный формат JSON")
		return
	}

	if err := payload.Validate(); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = storage.UpdateNote(h.DB, id, &payload)
	if err != nil {
		if err == sql.ErrNoRows {
			WriteError(w, http.StatusNotFound, "Заметка не найдена")
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
	id, err := parseID(w, r)
	if err != nil {
		return
	}

	err = storage.DeleteNote(h.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			WriteError(w, http.StatusNotFound, "Заметка не найдена")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Заметка с ID=%d успешно удалена", id),
	})
}

func parseID(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Неверный формат ID")
		return 0, err
	}

	return id, nil
}
