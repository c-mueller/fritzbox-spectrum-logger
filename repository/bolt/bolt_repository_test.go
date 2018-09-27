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
	"github.com/Flaque/filet"
	"github.com/c-mueller/fritzbox-spectrum-logger/fritz"
	"github.com/c-mueller/fritzbox-spectrum-logger/repository"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCollect_Stats(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	repo, err := NewBoltRepository(filepath.Join(tmpdir, "test_db.db"), false)
	assert.NoErrorf(t, err, "Initialization Failed")

	spectrum := loadTestSpectrum(t)

	insertSpectra(spectrum, repo, t, 10000, 6)

	stats, err := repo.GetStatistics()
	assert.NoError(t, err)

	assert.Equal(t, int64(10000), stats.TotalCount)

	repo.Close()
}

func TestInitRepo(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	repo, err := NewBoltRepository(filepath.Join(tmpdir, "test_db.db"), false)
	assert.NoErrorf(t, err, "Initialization Failed")
	repo.Close()
}

func TestRepository_Insert(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	repo, err := NewBoltRepository(filepath.Join(tmpdir, "test_db.db"), false)
	assert.NoErrorf(t, err, "Initialization Failed")

	spectrum := loadTestSpectrum(t)

	insertSpectra(spectrum, repo, t, 1000, 6)

	repo.Close()
}

func TestRepository_GetAllSpectrumKeys(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	repo, err := NewBoltRepository(filepath.Join(tmpdir, "test_db.db"), false)
	assert.NoErrorf(t, err, "Initialization Failed")

	spectrum := loadTestSpectrum(t)

	insertSpectra(spectrum, repo, t, 1000, 6)

	keys, _ := repo.GetAllSpectrumKeys()
	assert.True(t, len(keys) >= 250)
	for i, v := range keys {
		assert.Equal(t, "2018", v.Year)
		if i >= 1 {
			assert.True(t, keys.Less(i-1, i), "Array Not sorted Properly!")
		}
	}

	repo.Close()
}

func TestRepository_GetSpectrumForTimestamp(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	repo, err := NewBoltRepository(filepath.Join(tmpdir, "test_db.db"), false)
	assert.NoErrorf(t, err, "Initialization Failed")

	spectrum := loadTestSpectrum(t)

	err = repo.Insert(spectrum)
	assert.NoError(t, err)

	key, err := repo.GetSpectrum(1516233800)
	assert.NoError(t, err)
	assert.Equal(t, key.PortCount, 1)
	assert.Equal(t, key.Timestamp, int64(1516233800))

	repo.Close()
}

func TestRepository_GetSpectraForDayByTimestamp(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	repo, err := NewBoltRepository(filepath.Join(tmpdir, "test_db.db"), false)
	assert.NoErrorf(t, err, "Initialization Failed")

	spectrum := loadTestSpectrum(t)

	insertSpectra(spectrum, repo, t, 1000, 6)

	keys, _ := repo.GetAllSpectrumKeys()

	for _, v := range keys {
		spectraTimestamp, err := repo.GetTimestampsForSpectrumKey(v)
		assert.NoError(t, err)
		length := len(spectraTimestamp)
		if length == 0 {
			t.FailNow()
		}

		for _, ts := range spectraTimestamp {

			spec, err := repo.GetSpectrum(ts)
			assert.NoError(t, err)

			timestamp := time.Unix(spec.Timestamp, 0)
			y, m, d := v.GetIntegerValues()
			assert.Equal(t, y, timestamp.Year())
			assert.Equal(t, m, int(timestamp.Month()))
			assert.Equal(t, d, timestamp.Day())
		}
	}

	repo.Close()
}

func insertSpectra(spectrum *fritz.Spectrum, repo repository.Repository, t testing.TB, count, hourMultiplier int) {
	for i := 0; i < count; i++ {
		timestamp := time.Date(2018, 2, 14, 0, 0, 0, 0, time.UTC)
		timestamp = timestamp.Add(time.Duration(i*hourMultiplier) * time.Hour)
		spectrum.Timestamp = timestamp.Unix()
		err := repo.Insert(spectrum)
		assert.NoErrorf(t, err, "Inserting element %d failed", i)
	}
}

func loadTestSpectrum(t testing.TB) *fritz.Spectrum {
	file, err := os.Open("../testdata/example_spectrum.json")
	assert.NoError(t, err, "Loading Dummy Spectrum failed")
	var result *fritz.Spectrum
	data, err := ioutil.ReadAll(file)
	file.Close()
	err = json.Unmarshal(data, &result)
	return result
}
