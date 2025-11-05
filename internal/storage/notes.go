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

func GetNoteByID(db *sql.DB, id int) (models.Note, error) {
	var note models.Note
	err := db.QueryRow(`SELECT id, title, content, created_at FROM notes WHERE id=$1`, id).Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Заметка с id=%d не найдена", id)
			return note, err
		}
		log.Printf("Ошибка при получении заметки: %v", err)
		return note, err
	}
	return note, nil
}

func CreateNote(db *sql.DB, n *models.Note) (int, error) {
	var id int
	err := db.QueryRow("INSERT INTO notes (title, content) VALUES ($1, $2) RETURNING id", n.Title, n.Content).Scan(&id)
	if err != nil {
		log.Printf("Ошибка при создании заметки: %v", err)
		return 0, err
	}
	return id, nil
}

func UpdateNote(db *sql.DB, id int, n *models.Note) error {
	result, err := db.Exec("UPDATE notes SET title=$1, content=$2 WHERE id=$3", n.Title, n.Content, id)
	if err != nil {
		log.Printf("Ошибка при обновлении заметки id=%d", id)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("Заметка с id=%d не найдена", id)
		return sql.ErrNoRows
	}
	log.Printf("Успешно обновлена информация в заметке id=%d", id)
	return nil
}
