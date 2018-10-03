package fritz

import (
	"encoding/json"
	"fmt"
	"github.com/Flaque/filet"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestChartGenerator(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	files, err := ioutil.ReadDir("testdata/graphdata")
	assert.NoError(t, err)

	chartGen := NewConnectionSpeedChartGenerator("DSL Geschwindigkeitsverlauf", NewDefaultConnectionSpeedValueSet())

	for idx, f := range files {
		p := filepath.Join("testdata/graphdata", f.Name())
		fmt.Printf("%d/%d\n", idx+1, len(files))
		file, err := os.Open(p)
		assert.NoError(t, err)
		data, err := ioutil.ReadAll(file)
		assert.NoError(t, err)

		var spectrum Spectrum

		err = json.Unmarshal(data, &spectrum)

		chartGen.AddSpectrum(&spectrum)

		file.Close()
	}

	for _, v := range chartGen.ValueSets {
		fmt.Printf("%s: %f\n", v.Name, v.Values.AverageX())
	}

	path := filepath.Join(tmpdir, fmt.Sprintf("test-%d.png", time.Now().Unix()))

	w, err := os.Create(path)
	assert.NoError(t, err)

	data, err := chartGen.ToChart()
	assert.NoError(t, err)

	_, err = w.Write(data)
	assert.NoError(t, err)
	w.Close()
}
