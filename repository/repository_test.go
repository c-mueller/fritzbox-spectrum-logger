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
	"github.com/Flaque/filet"
	"github.com/c-mueller/fritzbox-spectrum-logger/fritz"
	"github.com/c-mueller/fritzbox-spectrum-logger/util"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInitRepo(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	repo, err := NewRepository(filepath.Join(tmpdir, "test_db.db"))
	assert.NoErrorf(t, err, "Initialization Failed")
	repo.Close()
}

func TestRepository_Insert(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	repo, err := NewRepository(filepath.Join(tmpdir, "test_db.db"))
	assert.NoErrorf(t, err, "Initialization Failed")

	spectrum := loadTestSpectrum(t)

	insertSpectra(spectrum, repo, t, 1000, 6)

	repo.Close()
}

func TestRepository_GetAllSpectrumKeys(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	repo, err := NewRepository(filepath.Join(tmpdir, "test_db.db"))
	assert.NoErrorf(t, err, "Initialization Failed")

	spectrum := loadTestSpectrum(t)

	insertSpectra(spectrum, repo, t, 1000, 6)

	keys, _ := repo.GetAllSpectrumKeys()
	assert.Equal(t, 251, len(keys))
	for i, v := range keys {
		assert.Equal(t, "2018", v.Year)
		if i >= 1 {
			assert.True(t, keys.Less(i-1, i), "Array Not sorted Properly!")
		}
	}

	repo.Close()
}

func TestRepository_GetSpectraForDay(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	repo, err := NewRepository(filepath.Join(tmpdir, "test_db.db"))
	assert.NoErrorf(t, err, "Initialization Failed")

	spectrum := loadTestSpectrum(t)

	insertSpectra(spectrum, repo, t, 1000, 6)

	for i := 1; i <= 5; i++ {
		s, err := repo.GetSpectraForDay(i, 3, 2018)
		assert.NoError(t, err, "Retrieving Spectras failed")
		assert.Equal(t, 4, len(s))

		for _, v := range s {
			assert.Equal(t, 1, v.PortCount)
			assert.NotEqual(t, 0, v.Ports[0].SpectrumInfo.PilotToneIndex)
		}
	}
	repo.Close()
}

func TestRepository_GetSpectraForDayByKeys(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	repo, err := NewRepository(filepath.Join(tmpdir, "test_db.db"))
	assert.NoErrorf(t, err, "Initialization Failed")

	spectrum := loadTestSpectrum(t)

	insertSpectra(spectrum, repo, t, 1000, 6)

	keys, _ := repo.GetAllSpectrumKeys()

	for _, v := range keys {
		spectra, err := repo.GetSpectraForSpectrumKey(v)
		assert.NoError(t, err)
		length := len(spectra)
		if length == 0 {
			t.FailNow()
		}

		for _, spec := range spectra {
			timestamp := time.Unix(spec.Timestamp, 0)
			y, m, d := v.GetIntegerValues()
			assert.Equal(t, y, timestamp.Year())
			assert.Equal(t, m, int(timestamp.Month()))
			assert.Equal(t, d, timestamp.Day())
		}
	}

	repo.Close()
}

func BenchmarkRepository_Insert(b *testing.B) {
	dir, err := ioutil.TempDir("", "")
	assert.NoError(b, err, "Creating tempdir failed")

	repo, err := NewRepository(filepath.Join(dir, "test_db.db"))
	assert.NoError(b, err, "Opening Repo Failed")

	spectrum := loadTestSpectrum(b)

	b.ResetTimer()

	count := 0
	start := time.Now()

	for i := 0; i < b.N; i++ {
		err = repo.Insert(spectrum)
		assert.NoError(b, err, "Insertion failed")
		spectrum.Timestamp = time.Now().Unix() + int64(100*i)
		count++
	}

	b.Log("Performed", count, "Operations During the benchmark")
	b.Log("The Benchmark ran", time.Since(start))

	repo.Close()
	err = util.RemoveContents(dir)
	if err != nil {
		b.Log("Cleanup failed!", err)
	}
}

func BenchmarkGetByKey(b *testing.B) {
	dir, err := ioutil.TempDir("", "")
	assert.NoError(b, err, "Creating tempdir failed")

	repo, err := NewRepository(filepath.Join(dir, "test_db.db"))
	assert.NoError(b, err, "Opening Repo Failed")

	spectrum := loadTestSpectrum(b)

	insertSpectra(spectrum, repo, b, b.N*5, 1)
	keys, _ := repo.GetAllSpectrumKeys()

	b.ResetTimer()

	start := time.Now()

	for i := 0; i < b.N; i++ {
		for _, v := range keys {
			spectra, err := repo.GetSpectraForSpectrumKey(v)
			assert.NoError(b, err)
			assert.True(b, len(spectra) > 0)
		}
	}

	b.Log("Element Count", b.N*5)
	b.Log("The Benchmark ran", time.Since(start))

	repo.Close()
	err = util.RemoveContents(dir)
	if err != nil {
		b.Log("Cleanup failed!", err)
	}
}

func insertSpectra(spectrum *fritz.Spectrum, repo *Repository, t testing.TB, count, hourMultiplier int) {
	for i := 0; i < count; i++ {
		timestamp := time.Now()
		timestamp = timestamp.Add(time.Duration(i*hourMultiplier) * time.Hour)
		spectrum.Timestamp = timestamp.Unix()
		err := repo.Insert(spectrum)
		assert.NoErrorf(t, err, "Inserting element %d failed", i)
	}
}

func loadTestSpectrum(t testing.TB) *fritz.Spectrum {
	file, err := os.Open("testdata/example_spectrum.json")
	assert.NoError(t, err, "Loading Dummy Spectrum failed")
	var result *fritz.Spectrum
	data, err := ioutil.ReadAll(file)
	file.Close()
	err = json.Unmarshal(data, &result)
	return result
}
