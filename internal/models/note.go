package models

import (
	"errors"
	"time"
)

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (n *Note) Validate() error {
	if n.Title == "" {
		return errors.New("Поле title обязательно")
	}
	if len(n.Title) > 100 {
		return errors.New("Поле title слишком длинное (максимум 100 символов)")
	}
	if n.Content == "" {
		return errors.New("Поле content обязательно")
	}

	return nil
}
