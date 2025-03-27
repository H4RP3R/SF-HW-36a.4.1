package memdb

import (
	"encoding/json"
	"fmt"
	"os"

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

	return posts, nil
}
