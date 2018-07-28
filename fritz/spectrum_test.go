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

package fritz

import (
	"encoding/json"
	"github.com/Flaque/filet"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const singlePortSpectrumPath = "testdata/example_spectrum.json"
const singlePortComparisonAPath = "testdata/example_spectrum_comparison_a.json"
const singlePortComparisonBPath = "testdata/example_spectrum_comparison_b.json"
const multiPortSpectrumPath = "testdata/example_spectrum_multiport.json"

const cleanup = false

func TestDrawComparison(t *testing.T) {
	testComparison(t, false, func(t *testing.T, data []byte) {
		assert.True(t, len(data) > 20000)
	})
}

func TestDrawComparison_Scaled(t *testing.T) {
	testComparison(t, true, func(t *testing.T, data []byte) {
		assert.True(t, len(data) > 50000)
	})
}

func TestDrawSpectrum(t *testing.T) {
	testRendering(t, cleanup, false, singlePortSpectrumPath, func(t *testing.T, data []byte) {
		assert.True(t, len(data) > 20000)
	})
}

func TestDrawSpectrum_Scaled(t *testing.T) {
	testRendering(t, cleanup, true, singlePortSpectrumPath, func(t *testing.T, data []byte) {
		assert.True(t, len(data) > 50000)
	})
}

func TestDrawSpectrum_MultiPort(t *testing.T) {
	testRendering(t, cleanup, false, multiPortSpectrumPath, func(t *testing.T, data []byte) {
		assert.True(t, len(data) > 40000)
	})
}
func TestDrawSpectrum_MultiPort_Scaled(t *testing.T) {
	testRendering(t, cleanup, true, multiPortSpectrumPath, func(t *testing.T, data []byte) {
		assert.True(t, len(data) > 100000)
	})
}

func BenchmarkRenderSpeed_SinglePort(b *testing.B) {
	benchmarkRendering(b, false, singlePortSpectrumPath)
}

func BenchmarkRenderSpeed_MultiPort(b *testing.B) {
	benchmarkRendering(b, false, multiPortSpectrumPath)
}

func BenchmarkRenderSpeed_SinglePort_Scaled(b *testing.B) {
	benchmarkRendering(b, true, singlePortSpectrumPath)
}

func BenchmarkRenderSpeed_MultiPort_Scaled(b *testing.B) {
	benchmarkRendering(b, true, multiPortSpectrumPath)
}

func testRendering(t *testing.T, cleanup, scaled bool, spectrumpath string, validator func(t *testing.T, data []byte)) {
	tmpDir := filet.TmpDir(t, "")
	t.Log("Using tmpdir", tmpDir)
	//Comment out the next line to Investigate the render Output
	if cleanup {
		defer filet.CleanUp(t)
	}

	t.Log("Loading test Data")
	data := loadTestData(t, spectrumpath)
	imgdata, err := data.Render(scaled)
	assert.NoError(t, err)
	t.Log("Data length:", len(imgdata), "Bytes")

	validator(t, imgdata)

	path := filepath.Join(tmpDir, "test.png")
	file, _ := os.Create(path)
	file.Write(imgdata)
	file.Close()
}

func benchmarkRendering(b *testing.B, scaled bool, path string) {
	b.Log("Loading test Data")
	data := loadTestData(nil, path)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data.Render(scaled)
	}
}

func loadTestData(t *testing.T, path string) *Spectrum {
	file, err := os.Open(path)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	var result *Spectrum
	data, err := ioutil.ReadAll(file)
	file.Close()
	err = json.Unmarshal(data, &result)
	return result
}

func testComparison(t *testing.T, scaled bool, validate func(t *testing.T, data []byte)) {
	tmpDir := filet.TmpDir(t, "")
	t.Log("Using tmpdir", tmpDir)
	//Comment out the next line to Investigate the render Output
	if cleanup {
		defer filet.CleanUp(t)
	}
	t.Log("Loading test Data")
	spectrumB := loadTestData(t, singlePortComparisonAPath)
	spectrumC := loadTestData(t, singlePortComparisonBPath)
	comparison := ComparisonSet{*spectrumB, *spectrumC}
	imgdata, err := comparison.RenderComparison(scaled)
	assert.NoError(t, err)
	validate(t, imgdata)
	path := filepath.Join(tmpDir, "test.png")
	file, _ := os.Create(path)
	file.Write(imgdata)
	file.Close()
}
