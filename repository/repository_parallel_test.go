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
	"fmt"
	"github.com/Flaque/filet"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

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
