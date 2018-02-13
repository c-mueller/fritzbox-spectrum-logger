package repository

import "github.com/boltdb/bolt"

type Repository struct {
	DatabasePath string
	db           *bolt.DB
}

type SpectrumKey struct {
	Year  string `json:"year"`
	Month string `json:"month"`
	Day   string `json:"day"`
}

type SpectraKeys []SpectrumKey
