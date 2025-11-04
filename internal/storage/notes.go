package storage

import (
	"database/sql"
	"log"

	"notes-api/internal/models"
)

func GetAllNotes(db *sql.DB) ([]models.Note, error) {
	rows, err := db.Query(
		`SELECT id, title, content, created_at
        FROM notes
        ORDER BY id; 
	`)
	if err != nil {
		log.Printf("Ошибка при запросе всех заметок: %v", err)
		return nil, err
	}

	defer rows.Close()
	var notes []models.Note

	for rows.Next() {
		var note models.Note
		err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)
		if err != nil {
			log.Printf("Ошибка при чтении строки: %v", err)
			return nil, err
		}
		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Ошибка при обходе строк: %v", err)
		return nil, err
	}

	return notes, nil
}
