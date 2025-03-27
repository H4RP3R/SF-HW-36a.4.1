package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"news/pkg/storage"
	"news/pkg/storage/memdb"
)

const testPostsPath = "../../test_data/post_examples.json"

func TestAPI_postsHandler(t *testing.T) {
	db := memdb.New()

	testPosts, err := memdb.LoadTestPosts(testPostsPath)
	if err != nil {
		t.Fatalf("unexpected error while loading test posts: %v", err)
	}

	err = db.AddPosts(testPosts)
	if err != nil {
		t.Fatalf("unexpected error while adding posts: %v", err)
	}

	api := New(db)
	path := fmt.Sprintf("/news/%d", len(testPosts))
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rr := httptest.NewRecorder()

	api.r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("want status code %v, got status code %v", http.StatusOK, rr.Code)
	}

	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("unexpected error while reading response body: %v", err)
	}

	var posts []storage.Post
	err = json.Unmarshal(b, &posts)
	if err != nil {
		t.Errorf("unexpected error while unmarshaling posts data: %v", err)
	}

	wantPosts := len(testPosts)
	gotPosts := len(posts)
	if wantPosts != gotPosts {
		t.Errorf("want %d posts, got %d posts", wantPosts, gotPosts)
	}
}
