package store

import (
	"sync"
	"time"

	"github.com/ASA11599/radius-server/internal/model"
)

type MemoryStore struct {
	posts []model.Post
	lock sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		posts: make([]model.Post, 0),
	}
}

func (ms *MemoryStore) SavePost(post model.Post) error {
	ms.deleteExpiredPosts()
	ms.lock.Lock()
	defer ms.lock.Unlock()
	ms.posts = append(ms.posts, post)
	return nil
}

func (ms *MemoryStore) deleteExpiredPosts() {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	activePosts := make([]model.Post, 0, len(ms.posts))
	for _, p := range ms.posts {
		now := time.Now().Unix()
		pDeadline := p.CreatedAt + p.Duration
		if now < pDeadline {
			activePosts = append(activePosts, p)
		}
	}
	ms.posts = activePosts
}

func (ms *MemoryStore) GetNearbyPosts(location model.Location, radius float64) ([]model.Post, error) {
	ms.deleteExpiredPosts()
	nearbyPosts := make([]model.Post, 0)
	for _, p := range ms.posts {
		if location.Distance(p.Location) <= radius {
			nearbyPosts = append(nearbyPosts, p)
		}
	}
	return nearbyPosts, nil
}

func (ms *MemoryStore) Ping() bool {
	return ms.posts != nil
}

func (ms *MemoryStore) Close() error {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	return nil
}
