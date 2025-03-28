package memdb

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"news/pkg/storage"
)

func LoadTestPosts(path string) ([]storage.Post, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to load test data from file: %w", err)
	}

	var posts []storage.Post
	err = json.Unmarshal(b, &posts)

	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal test data: %w", err)
	}

	for i := 0; i < len(posts); i++ {
		posts[i].Published = posts[i].Published.UTC()
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Published.After(posts[j].Published)
	})

	return posts, nil
}
