package handlers

import (
	"board/cache"
	"board/database"
	"board/models"
	"encoding/json"
	"net/http"
	"time"
)

var (
	threadCache         = cache.NewCache()
	threadAccessCounter = database.NewAccessCounter()
)

func GetThreads(w http.ResponseWriter, r *http.Request) {
	cacheKey := "threads"
	if cachedItem, found := threadCache.Get(cacheKey); found {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cachedItem)
		return
	}

	if threadAccessCounter.Increment(cacheKey) > 3 {
		rows, err := database.DB.Query("SELECT id, title, content, created_at FROM threads")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var threads []models.Thread
		for rows.Next() {
			var thread models.Thread
			var createdAt string
			if err := rows.Scan(&thread.ID, &thread.Title, &thread.Content, &createdAt); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			thread.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
			threads = append(threads, thread)
		}

		threadCache.Set(cacheKey, threads)
		threadAccessCounter.Reset(cacheKey)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(threads)
	} else {
		w.WriteHeader(http.StatusTooEarly)
	}
}

func CreateThread(w http.ResponseWriter, r *http.Request) {
	var thread models.Thread
	if err := json.NewDecoder(r.Body).Decode(&thread); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := database.DB.Prepare("INSERT INTO threads (title, content) VALUES (?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(thread.Title, thread.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	thread.ID = int(id)
	thread.CreatedAt = time.Now()

	// キャッシュをクリアする
	threadCache.Set("threads", nil)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(thread)
}
