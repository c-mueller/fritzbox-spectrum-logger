package repository

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/c-mueller/fritzbox-spectrum-logger/fritz"
	"github.com/op/go-logging"
	"time"
)

var log = logging.MustGetLogger("repository")

func NewRepository(path string) (*Repository, error) {
	log.Debugf("Opening database '%s'", path)
	db, err := bolt.Open(path, 0777, &bolt.Options{})

	if err != nil {
		return nil, err
	}

	err = initDb(db)

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) Insert(spectrum *fritz.Spectrum) error {
	tx, err := r.db.Begin(true)
	defer tx.Commit()
	if err != nil {
		return err
	}

	timestamp := time.Unix(spectrum.Timestamp, 0)
	year := fmt.Sprintf("%d", timestamp.Year())
	month := timestamp.Month().String()
	day := fmt.Sprintf("%d", timestamp.Day())

	spectraBucket, _ := tx.CreateBucketIfNotExists([]byte(SpectrumListBucketName))

	yearBucket, _ := spectraBucket.CreateBucketIfNotExists([]byte(year))
	monthBucket, _ := yearBucket.CreateBucketIfNotExists([]byte(month))
	dayBucket, _ := monthBucket.CreateBucketIfNotExists([]byte(day))

	jsonData, err := spectrum.JSON()
	if err != nil {
		return err
	}
	err = dayBucket.Put([]byte(string(fmt.Sprintf("%d", timestamp.Unix()))), jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Close() error {
	log.Debug("Closing Database")
	return r.db.Close()
}

func initDb(db *bolt.DB) error {
	tx, err := db.Begin(true)
	defer tx.Commit()

	if err != nil {
		return err
	}

	_, err = tx.CreateBucketIfNotExists([]byte(SpectrumListBucketName))
	if err != nil {
		return err
	}
	return nil
}
