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

package bolt

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/c-mueller/fritzbox-spectrum-logger/fritz"
	"github.com/c-mueller/fritzbox-spectrum-logger/repository"
	"github.com/c-mueller/fritzbox-spectrum-logger/repository/compress"
	"github.com/op/go-logging"
	"sort"
	"strconv"
	"time"
)

var log = logging.MustGetLogger("bolt_repository")

const (
	SpectrumListBucketName = "Spectra"
	SupportDataBucketName  = "SupportData"
)

type BoltRepository struct {
	DatabasePath string
	compress     bool
	db           *bolt.DB
}

func NewBoltRepository(path string, compress bool) (*BoltRepository, error) {
	log.Debugf("Opening database '%s'", path)
	db, err := bolt.Open(path, 0777, bolt.DefaultOptions)

	if err != nil {
		return nil, err
	}

	err = initDb(db)

	return &BoltRepository{
		db:       db,
		compress: compress,
	}, nil
}

func (r *BoltRepository) GetAllSpectrumKeys() (repository.SpectraKeys, error) {
	keys := make(repository.SpectraKeys, 0)

	err := r.forEachSpectrumKey(func(dayBucket *bolt.Bucket, key repository.SpectrumKey) error {
		keys = append(keys, key)
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Sort(keys)

	return keys, nil
}

func (r *BoltRepository) GetSpectrum(timestamp int64) (*fritz.Spectrum, error) {
	t := time.Unix(timestamp, 0)
	d, m, y := t.Day(), int(t.Month()), t.Year()
	return r.GetSpectrumWithDate(d, m, y, timestamp)
}

func (r *BoltRepository) GetSpectrumWithDate(day, month, year int, timestamp int64) (*fritz.Spectrum, error) {
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
			return repository.InvalidTimestampKey
		}

		if byteData[0] != []byte("{")[0] {
			byteData, err = compress.Decompress(byteData)
			if err != nil {
				return err
			}
		}

		err = json.Unmarshal(byteData, &spectrum)
		return err
	})
	if err != nil {
		return nil, err
	}

	return spectrum, nil
}

func (r *BoltRepository) GetTimestampsForDay(day, month, year int) (repository.TimestampArray, error) {
	data := make(repository.TimestampArray, 0)

	err := r.forEachSpectrumInDay(year, month, day, func(dayBucket *bolt.Bucket, k, v []byte) error {
		timestampString := string(k)
		timestamp, err := strconv.ParseInt(timestampString, 10, 64)
		if err != nil {
			return err
		}
		data = append(data, timestamp)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Sort(data)

	return data, nil
}

func (r *BoltRepository) Insert(spectrum *fritz.Spectrum) error {
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

		if r.compress {
			jsonData, err = compress.Compress(jsonData)
			if err != nil {
				return err
			}
		}

		err = dayBucket.Put([]byte(string(fmt.Sprintf("%d", timestamp.Unix()))), jsonData)
		return err
	})
	r.db.Sync()

	return err
}

func (r *BoltRepository) GetStatistics() (*repository.SpectraStats, error) {
	min := time.Now().Unix() * 2
	max := int64(0)
	count := int64(0)

	err := r.forEachSpectrumKey(func(dayBucket *bolt.Bucket, key repository.SpectrumKey) error {
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

	return &repository.SpectraStats{
		FirstSpectrum:  min,
		LatestSpectrum: max,
		TotalCount:     count,
	}, nil
}

func (r *BoltRepository) Close() error {
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

func (r *BoltRepository) StoreSupportData(data []byte, timestamp int) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(SupportDataBucketName))
		if err != nil {
			return err
		}

		cd, err := compress.Compress(data)

		if err != nil {
			return err
		}

		return b.Put([]byte(fmt.Sprintf("%d", timestamp)), cd)
	})
}

func (r *BoltRepository) ListSupportDataEntries() []int {
	data := make([]int, 0)

	r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(SupportDataBucketName))
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {

			d, _ := strconv.ParseInt(string(k), 10, 64)

			data = append(data, int(d))

			return nil
		})

	})

	return data
}

func (r *BoltRepository) GetSupportData(timestamp int) ([]byte, error) {
	var data []byte

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(SupportDataBucketName))
		if b == nil {
			return bolt.ErrBucketNotFound
		}

		data = b.Get([]byte(fmt.Sprintf("%d", timestamp)))

		return nil
	})

	if err != nil {
		return nil, err
	}

	return compress.Decompress(data)
}

func convertToByte(year, month, day int) ([]byte, []byte, []byte) {
	yearByte := []byte(fmt.Sprintf("%d", year))
	monthByte := []byte(fmt.Sprintf("%d", month))
	dayByte := []byte(fmt.Sprintf("%d", day))
	return yearByte, monthByte, dayByte
}
