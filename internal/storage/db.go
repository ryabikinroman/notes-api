package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func NewDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s post=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии соединения: %w", err)

	}

	for i := 0; i < 5; i++ {
		if err := db.Ping(); err != nil {
			log.Printf("БД недоступна (попытка %d/5): %v", i+1, err)
			time.Sleep(2 * time.Second)
		} else {
			log.Println("✅ Подключение к PostgreSQL установлено успешно")
			return db, nil
		}
	}
	return nil, fmt.Errorf("Не удалось подключиться к базе данных: %v", err)
}
