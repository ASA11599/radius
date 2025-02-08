package store

import "github.com/ASA11599/radius-server/internal/model"

type Store interface {
	Close() error
	SavePost(model.Post) error
	GetNearbyPosts(model.Location, float64) ([]model.Post, error)
	Ping() bool
}
