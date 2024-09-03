package models

type Note struct {
	Id      int64  `json:"id" db:"id"`
	UserId  int64  `json:"user_id" db:"user_id"`
	Content string `json:"content" db:"content"`
}
