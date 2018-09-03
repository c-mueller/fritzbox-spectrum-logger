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
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Benchmark_Bolt_Insert(b *testing.B) {
	dir := initializeTempDir(b)

	repo := initBoltRepository(dir, b)

	handleInsertionBenchmark(b, repo, dir)
}

func Benchmark_Bolt_Get_Timestamps(b *testing.B) {
	dir := initializeTempDir(b)

	repo := initBoltRepository(dir, b)

	handleBenchmarkGetTimestamps(b, repo, dir)
}

func Benchmark_Bolt_Retrieve(b *testing.B) {
	dir := initializeTempDir(b)

	repo := initBoltRepository(dir, b)

	handleRetrievalBenchmark(b, repo, dir)
}

func initBoltRepository(dir string, b *testing.B) *BoltRepository {
	repo, err := NewBoltRepository(filepath.Join(dir, "test_db.db"))
	assert.NoError(b, err, "Opening Repo Failed")
	return repo
}
