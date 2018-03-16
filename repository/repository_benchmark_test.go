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
	"github.com/c-mueller/fritzbox-spectrum-logger/util"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"testing"
	"time"
)

func Benchmark_Insert(b *testing.B) {
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

func Benchmark_Get_Timestamps(b *testing.B) {
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
			spectra, err := repo.GetTimestampsForSpectrumKey(v)
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

func Benchmark_Retrieve(b *testing.B) {
	dir, err := ioutil.TempDir("", "")
	assert.NoError(b, err, "Creating tempdir failed")

	repo, err := NewRepository(filepath.Join(dir, "test_db.db"))
	assert.NoError(b, err, "Opening Repo Failed")

	spectrum := loadTestSpectrum(b)

	insertSpectra(spectrum, repo, b, 5000, 1)
	keys, err := repo.GetAllSpectrumKeys()
	assert.NoError(b, err)

	timestamp, err := repo.GetTimestampsForSpectrumKey(keys[0])
	assert.NoError(b, err)

	rand.Seed(0xDEADBEEF)

	b.ResetTimer()

	start := time.Now()

	for i := 0; i < b.N; i++ {
		idx := rand.Intn(len(timestamp))
		spectra, err := repo.GetSpectrumForTimestamp(timestamp[idx])
		assert.NoError(b, err)
		assert.Equal(b, spectra.Timestamp, timestamp[idx])
		assert.Equal(b, spectra.PortCount, 1)
	}

	b.Log("Element Count", b.N*5)
	b.Log("The Benchmark ran", time.Since(start))

	repo.Close()
	err = util.RemoveContents(dir)
	if err != nil {
		b.Log("Cleanup failed!", err)
	}
}
