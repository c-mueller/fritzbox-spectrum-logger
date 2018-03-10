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
	"fmt"
	"github.com/Flaque/filet"
	"github.com/c-mueller/fritzbox-spectrum-logger/fritz"
	"github.com/c-mueller/fritzbox-spectrum-logger/util"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestCollect_Stats(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	repo, err := NewRepository(filepath.Join(tmpdir, "test_db.db"))
	assert.NoErrorf(t, err, "Initialization Failed")

	spectrum := loadTestSpectrum(t)

	insertSpectra(spectrum, repo, t, 10000, 6)

	stats, err := repo.GetStatistics()
	assert.NoError(t, err)

	assert.Equal(t, int64(10000), stats.TotalCount)

	repo.Close()
}

func Test_Parallel_Usage(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	repo, err := NewRepository(filepath.Join(tmpdir, "test_db.db"))
	assert.NoErrorf(t, err, "Initialization Failed")

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go handleInsertions(t, repo, &waitGroup)
	time.Sleep(100 * time.Millisecond)
	go handleRetrievals(t, repo, &waitGroup)

	waiter := make(chan struct{})

	go func() {
		defer close(waiter)
		waitGroup.Wait()
	}()
	select {
	case <-waiter:
		return
	case <-time.After(2 * time.Minute):
		t.Error("Test Timed out after 2 Minutes")
		t.FailNow()
	}

}

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
	assert.True(t, len(keys) >= 250)
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
		//assert.Equal(t, 4, len(s))

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

func handleInsertions(t *testing.T, repo *Repository, wg *sync.WaitGroup) {
	spectrum := loadTestSpectrum(t)
	defer wg.Done()
	cnt := 0
	for i := 0; i < 50; i++ {
		executions := int64(0)
		timeSum := int64(0)
		for j := 0; j < 250; j++ {
			spectrum.Timestamp = spectrum.Timestamp + int64(j+i*1000)

			ts := time.Now()
			err := repo.Insert(spectrum)
			timeSum += time.Since(ts).Nanoseconds()
			executions++

			assert.NoError(t, err)
			time.Sleep(10 * time.Microsecond)
			cnt++
		}

		fmt.Printf("Insertion Round %d Done. Average Operation Time %s\n", i, time.Duration(timeSum/executions)*time.Nanosecond)

		time.Sleep(100 * time.Microsecond)
	}
	t.Log("Inserted", cnt, "Elements")
}

func handleRetrievals(t *testing.T, repo *Repository, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 80; i++ {
		executions := int64(0)
		timeSum := int64(0)
		for j := 0; j < 500; j++ {
			ts := time.Now()
			_, err := repo.GetAllSpectrumKeys()
			timeSum += time.Since(ts).Nanoseconds()
			executions++
			assert.NoError(t, err)
		}
		fmt.Printf("Retrieval Round %d Done. Average Operation Time %s\n", i, time.Duration(timeSum/executions)*time.Nanosecond)
		time.Sleep(100 * time.Microsecond)
	}
}
