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

package relational

import (
	"encoding/json"
	"fmt"
	"github.com/c-mueller/fritzbox-spectrum-logger/fritz"
	"github.com/c-mueller/fritzbox-spectrum-logger/repository"
	"github.com/c-mueller/fritzbox-spectrum-logger/repository/compress"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("sql_repository")

type spectrumData struct {
	gorm.Model
	SpectrumData []byte `gorm:"size:20480"`
	Compressed   bool
}

type spectrumDSO struct {
	gorm.Model
	Year           int
	Month          int
	Day            int
	Timestamp      int64
	SpectrumDataID uint
}

type supportData struct {
	gorm.Model
	SupportData []byte `gorm:"size:8192000"`
	Timestamp   int    `gorm:"unique_index"`
}

type RelationalRepository struct {
	db       *gorm.DB
	compress bool
}

func (r *RelationalRepository) Backup() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.String(400, "Not Supported")
	}
}

func (r *RelationalRepository) StoreSupportData(data []byte, timestamp int) error {
	compressed, _ := compress.Compress(data)

	v := supportData{
		SupportData: compressed,
		Timestamp:   timestamp,
	}

	r.db.Save(&v)

	return nil
}

func (r *RelationalRepository) ListSupportDataEntries() []int {
	e := make([]supportData, 0)
	r.db.Find(&e)

	vals := make([]int, 0)
	for _, v := range e {
		vals = append(vals, v.Timestamp)
	}

	return vals
}

func (r *RelationalRepository) GetSupportData(timestamp int) ([]byte, error) {
	d := supportData{}
	r.db.Find(&d, supportData{Timestamp: timestamp})

	ucd, _ := compress.Decompress(d.SupportData)

	return ucd, nil
}

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
	db.AutoMigrate(&supportData{})
	db.Model(&spectrumDSO{}).AddUniqueIndex("idx_timestamp", "timestamp")

	log.Info("Initialized Database")
	return &RelationalRepository{
		db:       db,
		compress: compress,
	}, nil
}

func (r *RelationalRepository) GetAllSpectrumKeys() (repository.SpectraKeys, error) {
	keys := make([]spectrumDSO, 0)
	r.db.Find(&keys, &spectrumDSO{})

	sKeys := make(repository.SpectraKeys, 0)

	for _, v := range keys {
		sKeys = append(sKeys, repository.SpectrumKey{
			Year:  fmt.Sprintf("%d", v.Year),
			Month: fmt.Sprintf("%d", int(v.Month)),
			Day:   fmt.Sprintf("%d", v.Day),
		})
	}

	return sKeys, nil
}

func (r *RelationalRepository) GetTimestampsForSpectrumKey(key repository.SpectrumKey) (repository.TimestampArray, error) {
	y, m, d := key.GetIntegerValues()
	return r.GetTimestampsForDay(d, m, y)
}

func (r *RelationalRepository) GetSpectrum(timestamp int64) (*fritz.Spectrum, error) {
	var dso spectrumDSO

	r.db.First(&dso, &spectrumDSO{Timestamp: timestamp})

	return dso.toSpectrum(r)
}

func (r *RelationalRepository) GetTimestampsForDay(day, month, year int) (repository.TimestampArray, error) {
	data := make([]spectrumDSO, 0)
	r.db.Find(&data, &spectrumDSO{Day: day, Month: month, Year: year})

	timestamps := make([]int64, len(data))

	for k, v := range data {
		timestamps[k] = v.Timestamp
	}
	return repository.TimestampArray(timestamps), nil
}

func (r *RelationalRepository) Insert(spectrum *fritz.Spectrum) error {
	marshaledSpectrum, err := json.Marshal(spectrum)
	if err != nil {
		return err
	}

	data := marshaledSpectrum
	if r.compress {
		data, err = compress.Compress(data)
		if err != nil {
			return err
		}
	}

	spectrumData := spectrumData{
		SpectrumData: data,
		Compressed:   r.compress,
	}

	r.db.Create(&spectrumData)

	skey := repository.GetFromTimestamp(spectrum.Timestamp)
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

func (r *RelationalRepository) GetStatistics() (*repository.SpectraStats, error) {
	keys := make([]spectrumDSO, 0)
	r.db.Find(&keys, &spectrumDSO{})

	return &repository.SpectraStats{TotalCount: int64(len(keys)), FirstSpectrum: keys[0].Timestamp, LatestSpectrum: keys[len(keys)-1].Timestamp}, nil
}

func (r *RelationalRepository) Close() error {
	return r.db.Close()
}

func (dso *spectrumDSO) toSpectrum(repo *RelationalRepository) (*fritz.Spectrum, error) {
	var specData spectrumData
	repo.db.Find(&specData, dso.SpectrumDataID)

	data := specData.SpectrumData
	if specData.Compressed {
		dataI, err := compress.Decompress(data)
		if err != nil {
			return nil, err
		}
		data = dataI
	}

	var spectrum fritz.Spectrum
	err := json.Unmarshal(data, &spectrum)

	if err != nil {
		return nil, err
	}

	return &spectrum, nil
}
