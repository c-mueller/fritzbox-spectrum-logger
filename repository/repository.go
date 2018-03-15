// Fritz!Box Spectrum Logger (https://github.com/c-mueller/fritzbox-spectrum-logger).
// Copyright (c) 2018 Christian Müller <cmueller.dev@gmail.com>.
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
	"strconv"
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

	err := r.forEachSpectrumKey(func(dayBucket *bolt.Bucket, key SpectrumKey) error {
		keys = append(keys, key)
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Sort(keys)

	return keys, nil
}

func (r *Repository) GetSpectraForTimestamp(timestamp int64) (*fritz.Spectrum, error) {
	t := time.Unix(timestamp, 0)
	d, m, y := t.Day(), int(t.Month()), t.Year()
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
	data := make([]*fritz.Spectrum, 0)

	err := r.forEachSpectrumInDay(year, month, day, func(dayBucket *bolt.Bucket, k, v []byte) error {
		var spectrum *fritz.Spectrum
		err := json.Unmarshal(v, &spectrum)
		if err != nil {
			return err
		}
		data = append(data, spectrum)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *Repository) Insert(spectrum *fritz.Spectrum) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		timestamp := time.Unix(spectrum.Timestamp, 0)
		year, month, day := convertToByte(timestamp.Year(), int(timestamp.Month()), timestamp.Day())

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

func (r *Repository) GetStatistics() (*SpectraStats, error) {
	min := time.Now().Unix() * 2
	max := int64(0)
	count := int64(0)

	err := r.forEachSpectrumKey(func(dayBucket *bolt.Bucket, key SpectrumKey) error {
		err := dayBucket.ForEach(func(k, v []byte) error {
			count++
			parsedTimestamp, err := strconv.ParseInt(string(k), 10, 64)
			if err != nil {
				return err
			}
			if parsedTimestamp < min {
				min = parsedTimestamp
			}
			if parsedTimestamp > max {
				max = parsedTimestamp
			}
			return nil
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return &SpectraStats{
		FirstSpectrum:  min,
		LatestSpectrum: max,
		TotalCount:     count,
	}, nil
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
