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
	// AddPost adds a single post to the storage and returns the post ID and an error if any occurs.
	AddPost(post Post) (id uuid.UUID, err error)

	// AddPosts adds multiple posts to the storage and returns an error if any occurs.
	AddPosts(posts []Post) (err error)

	// Posts retrieves the 'n' newest posts from the storage and an error if any occurs.
	Posts(n int) (posts []Post, err error)
}
