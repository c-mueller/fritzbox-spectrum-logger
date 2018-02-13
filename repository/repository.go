package repository

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/c-mueller/fritzbox-spectrum-logger/fritz"
	"github.com/op/go-logging"
	"sort"
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

func (r *Repository) GetAllSpectrumKeys() (SpectraKeys, error) {
	tx, err := r.db.Begin(false)
	if err != nil {
		return nil, err
	}

	spectraBucket := tx.Bucket([]byte(SpectrumListBucketName))

	keys := make(SpectraKeys, 0)

	spectraBucket.ForEach(func(yearKey, v []byte) error {
		yearBucket := spectraBucket.Bucket(yearKey)
		if yearBucket == nil {
			//Ignore element if it is not a bucket
			return nil
		}
		yearBucket.ForEach(func(monthKey, v []byte) error {
			monthBucket := yearBucket.Bucket(monthKey)
			if monthBucket == nil {
				//Ignore element if it is not a bucket
				return nil
			}
			monthBucket.ForEach(func(dayKey, v []byte) error {
				dayBucket := monthBucket.Bucket(dayKey)
				if dayBucket == nil {
					//Ignore element if it is not a bucket
					return nil
				}
				key := SpectrumKey{
					Year:  string(yearKey),
					Month: string(monthKey),
					Day:   string(dayKey),
				}
				keys = append(keys, key)
				return nil
			})
			return nil
		})
		return nil
	})

	sort.Sort(keys)

	return keys, nil
}

func (r *Repository) GetSpectraForSpectrumKey(k SpectrumKey) ([]*fritz.Spectrum, error) {
	y, m, d := k.GetIntegerValues()
	return r.GetSpectraForDay(d, m, y)
}

func (r *Repository) GetSpectraForDay(day, month, year int) ([]*fritz.Spectrum, error) {
	yearByte := []byte(fmt.Sprintf("%d", year))
	monthByte := []byte(fmt.Sprintf("%d", month))
	dayByte := []byte(fmt.Sprintf("%d", day))

	tx, err := r.db.Begin(false)
	if err != nil {
		return nil, err
	}

	spectraBucket := tx.Bucket([]byte(SpectrumListBucketName))
	if spectraBucket == nil {
		log.Error("Spectra Bucket not found!")
		return nil, BucketNotFoundError
	}
	yearBucket := spectraBucket.Bucket(yearByte)
	if yearBucket == nil {
		log.Errorf("Year Bucket (Year: '%s') not found!",
			string(yearByte))
		return nil, BucketNotFoundError
	}
	monthBucket := yearBucket.Bucket(monthByte)
	if monthBucket == nil {
		log.Errorf("Month Bucket (Year: '%s' Month: '%s') not found!",
			string(yearByte), string(monthByte))
		return nil, BucketNotFoundError
	}
	dayBucket := monthBucket.Bucket(dayByte)
	if dayBucket == nil {
		log.Errorf("Month Bucket (Year: '%s' Month: '%s' Day: '%s') not found!",
			string(yearByte), string(monthByte), string(dayByte))
		return nil, BucketNotFoundError
	}

	data := make([]*fritz.Spectrum, 0)
	err = dayBucket.ForEach(func(k, v []byte) error {
		var spectrum *fritz.Spectrum
		err := json.Unmarshal(v, &spectrum)
		if err != nil {
			return err
		}
		data = append(data, spectrum)
		return nil
	})
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *Repository) Insert(spectrum *fritz.Spectrum) error {
	tx, err := r.db.Begin(true)
	defer tx.Commit()
	if err != nil {
		return err
	}

	timestamp := time.Unix(spectrum.Timestamp, 0)
	year := fmt.Sprintf("%d", timestamp.Year())
	month := fmt.Sprintf("%d", int(timestamp.Month()))
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
