package memdb

import (
	"sort"
	"sync"

	"news/pkg/storage"

	"github.com/gofrs/uuid"
)

type DB struct {
	mu    sync.Mutex
	posts map[uuid.UUID]storage.Post
}

func New() *DB {
	db := DB{
		posts: make(map[uuid.UUID]storage.Post),
	}

	return &db
}

func (db *DB) AddPost(post storage.Post) (id uuid.UUID, err error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	post.ID = uuid.NewV5(uuid.NamespaceURL, post.Link)
	db.posts[post.ID] = post

	return post.ID, nil
}

func (db *DB) AddPosts(posts []storage.Post) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for _, post := range posts {
		post.ID = uuid.NewV5(uuid.NamespaceURL, post.Link)
		db.posts[post.ID] = post
	}

	return
}

// Posts returns the n latest posts from the DB and an error if one occurs.
func (db *DB) Posts(n int) (posts []storage.Post, err error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for _, v := range db.posts {
		posts = append(posts, v)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Published.After(posts[j].Published)
	})

	if n > len(posts) {
		n = len(posts)
	}

	return posts[:n], nil
}
