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
	"github.com/c-mueller/fritzbox-spectrum-logger/fritz"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func NewRelationalRepository(mode, connectionString string, compress bool) (*RelationalRepository, error) {
	log.Debug("Connecting to SQL Database...")
	log.Debugf("Using mode %q with connection string %q", mode, connectionString)
	db, err := gorm.Open(mode, connectionString)
	if err != nil {
		return nil, err
	}

	log.Debug("Running Migrations (Creating Tables)...")
	db.AutoMigrate(&spectrumDSO{})
	db.Model(&spectrumDSO{}).AddUniqueIndex("idx_timestamp", "timestamp")

	log.Info("Initialized Database")
	return &RelationalRepository{
		db:       db,
		compress: compress,
	}, nil
}

func (r *RelationalRepository) GetAllSpectrumKeys() (SpectraKeys, error) {
	keys := make(SpectraKeys, 0)
	return keys, nil
}

func (r *RelationalRepository) GetSpectrumForTimestamp(timestamp int64) (*fritz.Spectrum, error) {
	var dso spectrumDSO

	r.db.First(&dso, &spectrumDSO{Timestamp: timestamp})

	return dso.toSpectrum()
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

	skey := GetFromTimestamp(spectrum.Timestamp)
	y, m, d := skey.GetIntegerValues()

	spectrumDSO := spectrumDSO{
		Timestamp:    spectrum.Timestamp,
		Day:          d,
		Month:        m,
		Year:         y,
		SpectrumData: data,
		Compressed:   r.compress,
	}

	r.db.Create(&spectrumDSO)

	return nil
}

func (r *RelationalRepository) GetStatistics() (*SpectraStats, error) {
	return &SpectraStats{TotalCount: 0, FirstSpectrum: 0, LatestSpectrum: 0}, nil
}

func (r *RelationalRepository) Close() error {
	return r.db.Close()
}
