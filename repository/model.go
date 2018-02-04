package repository

import "github.com/boltdb/bolt"

type Repository struct {
	DatabasePath string
	db           *bolt.DB
}

type RepositoryKey struct {
	Year  string `json:"year"`
	Month string `json:"month"`
	Day   string `json:"day"`
}
