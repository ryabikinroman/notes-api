package handlers

import (
	"encoding/json"
	"net/http"
	"notes-api/internal/storage"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Неверный JSON")
		return
	}

	if req.Email == "" || req.Password == "" {
		WriteError(w, http.StatusBadRequest, "Email и пароль обязательны")
		return
	}

	_, err := storage.GetUserByEmail(h.DB, req.Email)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Пользователь уже существует")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Ошибка хеширования пароля")
		return
	}

	userID, err := storage.CreateUser(h.DB, req.Email, string(hash))
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Ошибка создания пользователя")
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Регистрация успешна",
		"user_id": userID,
	})
}
