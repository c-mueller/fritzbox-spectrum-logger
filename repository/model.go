package repository

import "github.com/boltdb/bolt"

type Repository struct {
    DatabasePath string
    db           *bolt.DB
}
