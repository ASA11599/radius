package store

import (
	"math"
	"sync"
	"time"

	"github.com/ASA11599/radius-server/internal/model"
	"github.com/google/uuid"
)

type Index struct {
	locationMap map[[2]int][]uuid.UUID
}

func NewIndex() *Index {
	return &Index{
		locationMap: make(map[[2]int][]uuid.UUID),
	}
}

func (idx *Index) Add(post model.Post) {
	roughLat := int(math.Round(post.Location.Latitude))
	roughLon := int(math.Round(post.Location.Longitude))
	idx.locationMap[[2]int{ roughLat, roughLon }] = append(idx.locationMap[[2]int{roughLat, roughLon}], post.Id)
}

func (idx *Index) Delete(post model.Post) {
	roughLat := int(math.Round(post.Location.Latitude))
	roughLon := int(math.Round(post.Location.Longitude))
	ids := idx.locationMap[[2]int{ roughLat, roughLon }]
	for i, pid := range ids {
		if pid == post.Id {
			lastPostId := ids[len(ids) - 1]
			idx.locationMap[[2]int{ roughLat, roughLon }][i] = lastPostId
			idx.locationMap[[2]int{ roughLat, roughLon }] = ids[:len(ids) - 1]
			break
		}
	}
}

func (idx *Index) GetCandidates(location model.Location) []uuid.UUID {
	roughLat := int(math.Round(location.Latitude))
	roughLon := int(math.Round(location.Longitude))
	candidates := make([]uuid.UUID, 0)
	for dLat := -1; dLat < 2; dLat++ {
		for dLon := -1; dLon < 2; dLon++ {
			candidates = append(candidates, idx.locationMap[[2]int{ roughLat + dLat, roughLon + dLon }]...)
		}
	}
	return candidates
}

type IndexedMemoryStore struct {
	index *Index
	posts map[uuid.UUID]model.Post
	lock sync.Mutex
}

func NewIndexedMemoryStore() *IndexedMemoryStore {
	return &IndexedMemoryStore{
		index: NewIndex(),
		posts: make(map[uuid.UUID]model.Post),
	}
}

func (ims *IndexedMemoryStore) Ping() bool {
	return (ims.index != nil) && (ims.posts != nil)
}

func (ims *IndexedMemoryStore) SavePost(post model.Post) error {
	ims.deleteExpiredPosts()
	ims.lock.Lock()
	defer ims.lock.Unlock()
	ims.posts[post.Id] = post
	ims.index.Add(post)
	return nil
}

func (ims *IndexedMemoryStore) deleteExpiredPosts() {
	ims.lock.Lock()
	defer ims.lock.Unlock()
	for pid, p := range ims.posts {
		now := time.Now().Unix()
		pDeadline := p.CreatedAt + p.Duration
		if pDeadline < now {
			delete(ims.posts, pid)
			ims.index.Delete(p)
		}
	}
}

func (ims *IndexedMemoryStore) GetNearbyPosts(location model.Location, radius float64) ([]model.Post, error) {
	ims.deleteExpiredPosts()
	candidates := ims.index.GetCandidates(location)
	nearby := make([]model.Post, 0, len(candidates))
	for _, cid := range candidates {
		p := ims.posts[cid]
		if location.Distance(p.Location) <= radius {
			nearby = append(nearby, p)
		}
	}
	return nearby, nil
}

func (ims *IndexedMemoryStore) Close() error {
	ims.index = nil
	ims.posts = nil
	return nil
}
