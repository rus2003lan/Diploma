package report

import "sync"

type Repository struct {
	storage sync.Map
}

func New() *Repository {
	return &Repository{
		storage: sync.Map{},
	}
}