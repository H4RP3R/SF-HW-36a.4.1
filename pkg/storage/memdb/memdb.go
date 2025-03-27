package memdb

import (
	"sort"
	"sync"

	"news/pkg/storage"
)

type DB struct {
	mu    sync.Mutex
	posts map[int]storage.Post
	cntID int
}

func New() *DB {
	db := DB{
		posts: make(map[int]storage.Post),
		cntID: 1,
	}

	return &db
}

func (db *DB) AddPost(post storage.Post) (id int, err error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	id = db.cntID
	db.posts[id] = post
	db.cntID++

	return
}

func (db *DB) AddPosts(posts []storage.Post) (count int, err error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for _, post := range posts {
		db.posts[db.cntID] = post
		db.cntID++
		count++
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
		return posts[i].Published.Before(posts[j].Published)
	})

	if n > len(posts) {
		n = len(posts)
	}

	return posts[:n], nil
}
