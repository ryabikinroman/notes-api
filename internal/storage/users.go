package storage

import (
	"database/sql"
	"notes-api/internal/models"
)

func CreateUser(db *sql.DB, email, passwordHash string) (int, error) {
	var id int
	err := db.QueryRow(`
        INSERT INTO users (email, password_hash)
        VALUES ($1, $2)
        RETURNING id
    `, email, passwordHash).Scan(&id)

	return id, err
}

func GetUserByEmail(db *sql.DB, email string) (models.User, error) {
	var u models.User
	err := db.QueryRow(`
        SELECT id, email, password_hash
        FROM users WHERE email=$1
    `, email).Scan(&u.ID, &u.Email, &u.PasswordHash)

	return u, err
}
