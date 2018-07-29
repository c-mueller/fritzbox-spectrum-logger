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
const multiPortSpectrumPath = "testdata/example_spectrum_multiport.json"

func TestDrawSpectrum(t *testing.T) {
	tmpDir := filet.TmpDir(t, "")
	t.Log("Using tmpdir", tmpDir)
	//Comment out the next line to Investigate the render Output
	defer filet.CleanUp(t)

	t.Log("Loading test Data")
	data := loadTestData(t, singlePortSpectrumPath)
	imgdata, err := data.Render()
	assert.NoError(t, err)
	t.Log("Data length:", len(imgdata), "Bytes")

	assert.True(t, len(imgdata) > 50000)

	path := filepath.Join(tmpDir, "test.png")
	file, _ := os.Create(path)
	file.Write(imgdata)
	file.Close()
}

func TestDrawSpectrum_MultiPort(t *testing.T) {
	tmpDir := filet.TmpDir(t, "")
	t.Log("Using tmpdir", tmpDir)
	//Comment out the next line to Investigate the render Output
	defer filet.CleanUp(t)

	t.Log("Loading test Data")
	data := loadTestData(t, multiPortSpectrumPath)
	imgdata, err := data.Render()
	assert.NoError(t, err)
	t.Log("Data length:", len(imgdata), "Bytes")

	assert.True(t, len(imgdata) > 50000)

	path := filepath.Join(tmpDir, "test.png")
	file, _ := os.Create(path)
	file.Write(imgdata)
	file.Close()
}

func BenchmarkRenderSpeed_SinglePort(b *testing.B) {
	b.Log("Loading test Data")
	data := loadTestData(nil, singlePortSpectrumPath)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data.Render()
	}
}

func BenchmarkRenderSpeed_MultiPort(b *testing.B) {
	b.Log("Loading test Data")
	data := loadTestData(nil, multiPortSpectrumPath)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data.Render()
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
