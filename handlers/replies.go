package handlers

import (
	"board/cache"
	"board/database"
	"board/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

var (
	replyCache         = cache.NewCache()
	replyAccessCounter = database.NewAccessCounter()
)

func GetReplies(w http.ResponseWriter, r *http.Request) {
	threadID := r.URL.Query().Get("thread_id")
	cacheKey := "replies_" + threadID
	if cachedItem, found := replyCache.Get(cacheKey); found {
		println("this is cache data")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cachedItem)
		return
	}

	if replyAccessCounter.Increment(cacheKey) < 3 {
		rows, err := database.DB.Query("SELECT id, thread_id, content, created_at FROM replies WHERE thread_id = ?", threadID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var replies []models.Reply
		for rows.Next() {
			var reply models.Reply
			var createdAt string
			if err := rows.Scan(&reply.ID, &reply.ThreadID, &reply.Content, &createdAt); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			reply.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
			replies = append(replies, reply)
		}

		replyCache.Set(cacheKey, replies)
		replyAccessCounter.Reset(cacheKey)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(replies)
	} else {
		w.WriteHeader(http.StatusTooEarly)
	}
}

func CreateReply(w http.ResponseWriter, r *http.Request) {
	var reply models.Reply
	if err := json.NewDecoder(r.Body).Decode(&reply); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var maxID int
	err := database.DB.QueryRow("SELECT IFNULL(MAX(id), 0) + 1 FROM replies WHERE thread_id = ?", reply.ThreadID).Scan(&maxID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reply.ID = maxID

	stmt, err := database.DB.Prepare("INSERT INTO replies (id, thread_id, content) VALUES (?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(reply.ID, reply.ThreadID, reply.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reply.CreatedAt = time.Now()

	// キャッシュをクリアする
	replyCache.Set("replies_"+strconv.Itoa(reply.ThreadID), nil)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reply)
}
