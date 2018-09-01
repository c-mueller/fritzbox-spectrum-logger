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
	"github.com/c-mueller/fritzbox-spectrum-logger/fritz"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

func NewRelationalRepository(path string) (*RelationalRepository, error) {
	gorm.Open("sqlite", "test.db")
	return nil, nil
}

func (r *RelationalRepository) GetAllSpectrumKeys() (SpectraKeys, error) {
	keys := make(SpectraKeys, 0)
	return keys, nil
}

func (r *RelationalRepository) GetSpectrumForTimestamp(timestamp int64) (*fritz.Spectrum, error) {
	t := time.Unix(timestamp, 0)
	d, m, y := t.Day(), int(t.Month()), t.Year()
	return r.GetSpectrum(d, m, y, timestamp)
}

func (r *RelationalRepository) GetSpectrum(day, month, year int, timestamp int64) (*fritz.Spectrum, error) {
	return nil, nil
}

func (r *RelationalRepository) GetTimestampsForDay(day, month, year int) (TimestampArray, error) {
	return nil, nil
}

func (r *RelationalRepository) Insert(spectrum *fritz.Spectrum) error {
	return nil
}

func (r *RelationalRepository) GetStatistics() (*SpectraStats, error) {

	return nil, nil
}

func (r *RelationalRepository) Close() error {

	return nil
}
