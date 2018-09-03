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
	"github.com/c-mueller/fritzbox-spectrum-logger/fritz"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func NewSQLiteRepository(path string, compress bool) (*RelationalRepository, error) {
	return NewRelationalRepository("sqlite3", path, compress)
}

func NewRelationalRepository(mode, connectionString string, compress bool) (*RelationalRepository, error) {
	log.Debug("Connecting to SQL Database...")
	log.Debugf("Using mode %q with connection string %q", mode, connectionString)
	db, err := gorm.Open(mode, connectionString)
	if err != nil {
		return nil, err
	}

	log.Debug("Running Migrations (Creating Tables)...")
	db.AutoMigrate(&spectrumDSO{})
	db.AutoMigrate(&spectrumData{})
	db.Model(&spectrumDSO{}).AddUniqueIndex("idx_timestamp", "timestamp")

	log.Info("Initialized Database")
	return &RelationalRepository{
		db:       db,
		compress: compress,
	}, nil
}

func (r *RelationalRepository) GetAllSpectrumKeys() (SpectraKeys, error) {
	keys := make([]spectrumDSO, 0)
	r.db.Find(&keys, &spectrumDSO{})

	sKeys := make(SpectraKeys, 0)

	for _, v := range keys {
		sKeys = append(sKeys, SpectrumKey{
			Year:  fmt.Sprintf("%d", v.Year),
			Month: fmt.Sprintf("%d", int(v.Month)),
			Day:   fmt.Sprintf("%d", v.Day),
		})
	}

	return sKeys, nil
}

func (r *RelationalRepository) GetTimestampsForSpectrumKey(key SpectrumKey) (TimestampArray, error) {
	y, m, d := key.GetIntegerValues()
	return r.GetTimestampsForDay(d, m, y)
}

func (r *RelationalRepository) GetSpectrumForTimestamp(timestamp int64) (*fritz.Spectrum, error) {
	var dso spectrumDSO

	r.db.First(&dso, &spectrumDSO{Timestamp: timestamp})

	return dso.toSpectrum(r)
}

func (r *RelationalRepository) GetSpectrum(day, month, year int, timestamp int64) (*fritz.Spectrum, error) {
	return r.GetSpectrumForTimestamp(timestamp)
}

func (r *RelationalRepository) GetTimestampsForDay(day, month, year int) (TimestampArray, error) {
	data := make([]spectrumDSO, 0)
	r.db.Find(&data, &spectrumDSO{Day: day, Month: month, Year: year})

	timestamps := make([]int64, len(data))

	for k, v := range data {
		timestamps[k] = v.Timestamp
	}
	return TimestampArray(timestamps), nil
}

func (r *RelationalRepository) Insert(spectrum *fritz.Spectrum) error {
	marshaledSpectrum, err := json.Marshal(spectrum)
	if err != nil {
		return err
	}

	data := marshaledSpectrum
	if r.compress {
		data, err = compress(data)
		if err != nil {
			return err
		}
	}

	spectrumData := spectrumData{
		SpectrumData: data,
		Compressed:   r.compress,
	}

	r.db.Create(&spectrumData)

	skey := GetFromTimestamp(spectrum.Timestamp)
	y, m, d := skey.GetIntegerValues()

	spectrumDSO := spectrumDSO{
		Timestamp:      spectrum.Timestamp,
		Day:            d,
		Month:          m,
		Year:           y,
		SpectrumDataID: spectrumData.ID,
	}

	r.db.Create(&spectrumDSO)

	return nil
}

func (r *RelationalRepository) GetStatistics() (*SpectraStats, error) {
	keys := make([]spectrumDSO, 0)
	r.db.Find(&keys, &spectrumDSO{})

	return &SpectraStats{TotalCount: int64(len(keys)), FirstSpectrum: keys[0].Timestamp, LatestSpectrum: keys[len(keys)-1].Timestamp}, nil
}

func (r *RelationalRepository) Close() error {
	return r.db.Close()
}
