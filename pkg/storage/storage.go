package storage

import "time"

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Published time.Time `json:"published"`
	Link      string    `json:"link"`
}

type Storage interface {
	AddPost(post Post) (id int, err error)
	AddPosts(posts []Post) (count int, err error)
	Posts(n int) (posts []Post, err error)
}
