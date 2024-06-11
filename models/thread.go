package models

import "time"

type Thread struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Reply struct {
	ID        int       `json:"id"`
	ThreadID  int       `json:"thread_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
