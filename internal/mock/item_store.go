package mock

import (
	"location-service/internal/types"
)

type ItemStore struct {
	store                   map[string]string
	AddNewFn                func() error
	GetFn                   func() (error)
	GetAllNearbyFn          func() ([]string, error)
	GetAllNearbyUnmatchedFn func() ([]string, error)
	UpdateFn                func() error
	DeleteFn                func() error
	SetMatchedFn            func() error
	SetUnmatchedFn          func() error
}

func NewItemStore() *ItemStore {
	return &ItemStore{
		store: make(map[string]string),
	}
}

func (m *ItemStore) AddNew(t types.TrackedItem) error {
	return m.AddNewFn()
}

func (m *ItemStore) Get(id string, t types.TrackedItem) error {
	return m.GetFn()
}

func (m *ItemStore) GetAllNearby(coord map[string]float64, radius float64) ([]string, error) {
	return m.GetAllNearbyFn()
}

func (m *ItemStore) GetAllNearbyUnmatched(coord map[string]float64, radius float64) ([]string, error) {
	return m.GetAllNearbyUnmatchedFn()
}

func (m *ItemStore) Update(t types.TrackedItem) error {
	return m.UpdateFn()
}

func (m *ItemStore) Delete(id string) error {
	return m.DeleteFn()
}

// handle setting matched twice
func (m *ItemStore) SetMatched(id string) error {
	return m.SetMatchedFn()
}

// handle setting unmatched when already unmatched, and when matched...
func (m *ItemStore) SetUnmatched(id string) error {
	return m.SetUnmatchedFn()
}
