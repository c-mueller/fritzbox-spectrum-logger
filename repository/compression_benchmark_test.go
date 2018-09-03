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
	"encoding/json"
	"github.com/c-mueller/fritzbox-spectrum-logger/fritz"
	"testing"
	"time"
)

func Benchmark_CompressionSpeed(b *testing.B) {
	spectrum := loadTestSpectrum(b)

	start := time.Now()
	opcnt := 0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data, _ := json.Marshal(spectrum)
		compress(data)

		opcnt++
	}

	execTime := time.Now().Sub(start)

	b.Logf("Ran %d compression Iterations in %s", opcnt, execTime.String())
}

func Benchmark_DecompressionSpeed(b *testing.B) {
	spectrum := loadTestSpectrum(b)
	data, _ := json.Marshal(spectrum)
	compressedData, _ := compress(data)

	start := time.Now()
	opcnt := 0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uncompressed, _ := decompress(compressedData)
		var spectrum fritz.Spectrum
		json.Unmarshal(uncompressed, &spectrum)

		opcnt++
	}

	execTime := time.Now().Sub(start)

	b.Logf("Ran %d decompression Iterations in %s", opcnt, execTime.String())
}
