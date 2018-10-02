package fritz

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestChartGenerator(t *testing.T) {
	files, err := ioutil.ReadDir("testdata/graphdata")
	assert.NoError(t, err)

	chartGen := ChartGenerator{
		Title: "DSL Geschwindigkeitsverlauf - August bis September 2018",
		ValueSets: []*ValueSet{
			{
				Name: "Aktuelle Datenrate - Downstream",
				Extractor: func(s *Spectrum, conInfo *ConnectionInformation) (float64, error) {
					return float64(conInfo.Downstream.CurrentDataRate), nil
				},
			},
			{
				Name: "Aktuelle Datenrate - Upstream",
				Extractor: func(s *Spectrum, conInfo *ConnectionInformation) (float64, error) {
					return float64(conInfo.Upstream.CurrentDataRate), nil
				},
			},
			{
				Name: "Leitungskapazität - Downstream",
				Extractor: func(s *Spectrum, conInfo *ConnectionInformation) (float64, error) {
					return float64(conInfo.Downstream.Capacity), nil
				},
			},
			{
				Name: "Leitungskapazität - Upstream",
				Extractor: func(s *Spectrum, conInfo *ConnectionInformation) (float64, error) {
					return float64(conInfo.Upstream.Capacity), nil
				},
			},
			//{
			//	Name: "Maximale Datenrate - Downstream",
			//	Extractor: func(s *Spectrum, conInfo *ConnectionInformation) (float64, error) {
			//		return float64(conInfo.Downstream.MaximumDataRate), nil
			//	},
			//},
			//{
			//	Name: "Maximale Datenrate - Upstream",
			//	Extractor: func(s *Spectrum, conInfo *ConnectionInformation) (float64, error) {
			//		return float64(conInfo.Upstream.MaximumDataRate), nil
			//	},
			//},
		},
	}

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

	w, err := os.Create(fmt.Sprintf("test-%d.png", time.Now().Unix()))
	assert.NoError(t, err)

	data, err := chartGen.ToChart()
	assert.NoError(t, err)

	_, err = w.Write(data)
	assert.NoError(t, err)
	w.Close()
}
