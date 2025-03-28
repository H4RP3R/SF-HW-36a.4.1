package storage

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

var (
	ErrConnectDB       = fmt.Errorf("unable to establish DB connection")
	ErrDBNotResponding = fmt.Errorf("DB not responding")
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Published time.Time `json:"published"`
	Link      string    `json:"link"`
}

type Storage interface {
	AddPost(post Post) (id uuid.UUID, err error)
	AddPosts(posts []Post) (err error)
	Posts(n int) (posts []Post, err error)
}
