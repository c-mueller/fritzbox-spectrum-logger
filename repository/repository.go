// Fritz!Box Spectrum Logger (https://github.com/c-mueller/fritzbox-spectrum-logger).
// Copyright (c) 2018 Christian Müller<cmueller.dev@gmail.com>.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

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
	db, err := bolt.Open(path, 0777, bolt.DefaultOptions)

	if err != nil {
		return nil, err
	}

	err = initDb(db)

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetAllSpectrumKeys() (SpectraKeys, error) {
	keys := make(SpectraKeys, 0)

	err := r.db.View(func(tx *bolt.Tx) error {
		spectraBucket := tx.Bucket([]byte(SpectrumListBucketName))

		err := spectraBucket.ForEach(func(yearKey, v []byte) error {
			yearBucket := spectraBucket.Bucket(yearKey)
			if yearBucket == nil {
				//Ignore element if it is not a bucket
				return nil
			}
			err := yearBucket.ForEach(func(monthKey, v []byte) error {
				monthBucket := yearBucket.Bucket(monthKey)
				if monthBucket == nil {
					//Ignore element if it is not a bucket
					return nil
				}
				err := monthBucket.ForEach(func(dayKey, v []byte) error {
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
				return err
			})
			return err
		})
		return err
	})

	if err != nil {
		return nil, err
	}

	sort.Sort(keys)

	return keys, nil
}

func (r *Repository) GetSpectraForSpectrumKey(k SpectrumKey) ([]*fritz.Spectrum, error) {
	if !k.IsValid() {
		return nil, InvalidDateKey
	}
	y, m, d := k.GetIntegerValues()
	return r.GetSpectraForDay(d, m, y)
}

func (r *Repository) GetSpectrumBySpectrumKey(k *SpectrumKey, timestamp int64) (*fritz.Spectrum, error) {
	if !k.IsValid() {
		return nil, InvalidDateKey
	}
	y, m, d := k.GetIntegerValues()
	return r.GetSpectrum(d, m, y, timestamp)
}

func (r *Repository) GetSpectrum(day, month, year int, timestamp int64) (*fritz.Spectrum, error) {
	yearByte, monthByte, dayByte := convertToByte(year, month, day)
	timestampByte := []byte(fmt.Sprintf("%d", timestamp))
	var spectrum *fritz.Spectrum

	err := r.db.View(func(tx *bolt.Tx) error {
		dayBucket, err := r.getDayBucket(dayByte, monthByte, yearByte, tx)
		if err != nil {
			return err
		}
		byteData := dayBucket.Get(timestampByte)
		if byteData == nil {
			return InvalidTimestampKey
		}
		err = json.Unmarshal(byteData, &spectrum)
		return err
	})
	if err != nil {
		return nil, err
	}

	return spectrum, nil
}

func (r *Repository) GetSpectraForDay(day, month, year int) ([]*fritz.Spectrum, error) {
	yearByte, monthByte, dayByte := convertToByte(year, month, day)
	data := make([]*fritz.Spectrum, 0)

	err := r.db.View(func(tx *bolt.Tx) error {

		dayBucket, err := r.getDayBucket(dayByte, monthByte, yearByte, tx)
		if err != nil {
			return err
		}

		err = dayBucket.ForEach(func(k, v []byte) error {
			var spectrum *fritz.Spectrum
			err := json.Unmarshal(v, &spectrum)
			if err != nil {
				return err
			}
			data = append(data, spectrum)
			return nil
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *Repository) Insert(spectrum *fritz.Spectrum) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
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
		return err
	})
	r.db.Sync()

	return err
}

func (r *Repository) Close() error {
	log.Debug("Closing Database")
	return r.db.Close()
}

func (r *Repository) getDayBucket(dayByte, monthByte, yearByte []byte, tx *bolt.Tx) (*bolt.Bucket, error) {
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
	return dayBucket, nil
}

func convertToByte(year, month, day int) ([]byte, []byte, []byte) {
	yearByte := []byte(fmt.Sprintf("%d", year))
	monthByte := []byte(fmt.Sprintf("%d", month))
	dayByte := []byte(fmt.Sprintf("%d", day))
	return yearByte, monthByte, dayByte
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
