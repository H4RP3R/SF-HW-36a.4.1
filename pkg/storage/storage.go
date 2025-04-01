package storage

import (
	"fmt"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"

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

// ValidatePosts accepts a slice of posts and removes the invalid ones, i.e., posts containing any empty fields.
func ValidatePosts(posts ...Post) []Post {
	var validPosts []Post
	for _, p := range posts {
		if p.Title != "" && p.Content != "" && p.Link != "" && !p.Published.IsZero() {
			if _, err := url.ParseRequestURI(p.Link); err == nil {
				validPosts = append(validPosts, p)
			}
		} else {
			log.Warnf("Invalidated post: %+v", p)
		}
	}

	return validPosts
}
