// Fritz!Box Spectrum Logger (https://github.com/c-mueller/fritzbox-spectrum-logger).
// Copyright (c) 2018 Christian Müller <cmueller.dev@gmail.com>.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package repository

import (
	"fmt"
	"github.com/Flaque/filet"
	"github.com/c-mueller/fritzbox-spectrum-logger/fritz"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSQLiteInsert_Compressed(t *testing.T) {
	//defer filet.CleanUp(t)

	testSpectrum := loadTestSpectrum(t)

	db := initializeSQLiteDatabase(t, true)
	defer db.Close()

	handleSQLInsertions(t, testSpectrum, db)

	validateInsertions(t, db, 100, 1000, 1500)
}

func TestSQLiteInsert_Uncompressed(t *testing.T) {
	//defer filet.CleanUp(t)

	testSpectrum := loadTestSpectrum(t)

	db := initializeSQLiteDatabase(t, false)
	defer db.Close()

	handleSQLInsertions(t, testSpectrum, db)

	validateInsertions(t, db, 100, 8000, 9500)
}

func TestSQLiteQueryByTimestamp_Compressed(t *testing.T) {
	//defer filet.CleanUp(t)

	testSpectrum := loadTestSpectrum(t)

	db := initializeSQLiteDatabase(t, true)
	defer db.Close()

	err := db.Insert(testSpectrum)
	assert.NoError(t, err)

	spc, err := db.GetSpectrumForTimestamp(testSpectrum.Timestamp)
	assert.NoError(t, err)

	assert.Equal(t, testSpectrum.Timestamp, spc.Timestamp)
	assert.Equal(t, testSpectrum.ConnectionInformation, spc.ConnectionInformation)
}

func TestSQLiteQueryByTimestamp_Uncompressed(t *testing.T) {
	//defer filet.CleanUp(t)

	testSpectrum := loadTestSpectrum(t)

	db := initializeSQLiteDatabase(t, false)
	defer db.Close()

	err := db.Insert(testSpectrum)
	assert.NoError(t, err)

	spc, err := db.GetSpectrumForTimestamp(testSpectrum.Timestamp)
	assert.NoError(t, err)

	assert.Equal(t, testSpectrum.Timestamp, spc.Timestamp)
	assert.Equal(t, testSpectrum.ConnectionInformation, spc.ConnectionInformation)
}

func handleSQLInsertions(t *testing.T, testSpectrum *fritz.Spectrum, db *RelationalRepository) {
	t.Log("Inserting 100 Spectra")

	for i := 0; i < 100; i++ {
		testSpectrum.Timestamp = testSpectrum.Timestamp + int64(10*i)

		err := db.Insert(testSpectrum)
		assert.NoError(t, err)
	}
}

func validateInsertions(t *testing.T, repo *RelationalRepository, count, dataSizeMin, dataSizeMax int) {
	t.Log("Validating Insertions")

	spectra := make([]spectrumDSO, 0)

	repo.db.Find(&spectra, &spectrumDSO{})

	assert.Equal(t, count, len(spectra))

	for _, v := range spectra {
		var d spectrumData
		repo.db.Find(&d, v.SpectrumDataID)
		assert.True(t, len(d.SpectrumData) >= dataSizeMin && len(d.SpectrumData) <= dataSizeMax,
			"Data Size out of bounds expected %d <= x <= %d. But x is %d", dataSizeMin, dataSizeMax, len(d.SpectrumData))
	}

}

func initializeSQLiteDatabase(t *testing.T, compress bool) *RelationalRepository {
	tmpdir := filet.TmpDir(t, "")
	t.Log("Using tempdir", tmpdir)
	databasePath := fmt.Sprintf("%s/test_database.db", tmpdir)

	db, err := NewSQLiteRepository(databasePath, compress)
	if err != nil {
		assert.NoError(t, err)
		t.FailNow()
	}
	return db
}
