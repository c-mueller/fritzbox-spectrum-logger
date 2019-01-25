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
	"github.com/gin-gonic/gin"
)

type RepoBuilder interface {
	BuildRepository(compress bool, path string) (Repository, error)
	GetName() string
}

type Repository interface {
	GetAllSpectrumKeys() (SpectraKeys, error)
	GetSpectrum(timestamp int64) (*fritz.Spectrum, error)
	GetTimestampsForDay(day, month, year int) (TimestampArray, error)
	GetTimestampsForSpectrumKey(key SpectrumKey) (TimestampArray, error)
	Insert(spectrum *fritz.Spectrum) error
	GetStatistics() (*SpectraStats, error)

	StoreSupportData(data []byte, timestamp int) error
	ListSupportDataEntries() []int
	GetSupportData(timestamp int) ([]byte, error)

	Backup() gin.HandlerFunc

	Close() error
}
