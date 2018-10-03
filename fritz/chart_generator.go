package fritz

import (
	"bytes"
	"github.com/wcharczuk/go-chart"
	"time"
)

type ChartGenerator struct {
	Width     int
	Height    int
	XAxisName string
	YAxisName string
	ValueSets []*ValueSet
	Title     string
}

type ValueSet struct {
	Name      string
	Extractor func(s *Spectrum, conInfo *ConnectionInformation) (float64, error)
	Values    ValuePairs
}

type ValuePairs []ValuePair

type ValuePair struct {
	Timestamp int64
	Value     float64
}

func NewDefaultConnectionSpeedValueSet() []*ValueSet {
	return []*ValueSet{
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
	}
}

func NewConnectionSpeedChartGenerator(title string, valueSets []*ValueSet) *ChartGenerator {
	return &ChartGenerator{
		XAxisName: "Zeit",
		YAxisName: "KBit/s",
		Width:     5000,
		Height:    2000,
		ValueSets: valueSets,
	}
}

func (v ValuePairs) AverageX() float64 {
	average := float64(0)
	for _, v := range v {
		average += v.Value
	}

	average = average / float64(len(v))

	return average
}

func (v ValuePairs) toAxisValues() ([]time.Time, []float64) {
	x := make([]time.Time, len(v))
	y := make([]float64, len(v))

	for k, v := range v {
		x[k] = time.Unix(v.Timestamp, 0)
		y[k] = v.Value
	}

	return x, y
}

func (c *ChartGenerator) AddSpectrum(s *Spectrum) {
	conInfo, err := s.GetConnectionInformation()
	if err != nil {
		return
	}
	for _, valueSet := range c.ValueSets {
		value, err := valueSet.Extractor(s, conInfo)
		if err != nil {
			value = -1
		}

		if valueSet.Values == nil {
			valueSet.Values = make(ValuePairs, 0)
		}

		valueSet.Values = append(valueSet.Values, ValuePair{s.Timestamp, value})
	}
}

func (c *ChartGenerator) ToChart() ([]byte, error) {

	series := make([]chart.Series, 0)

	for _, set := range c.ValueSets {
		x, y := set.Values.toAxisValues()

		s := chart.TimeSeries{
			Name:    set.Name,
			XValues: x,
			YValues: y,
		}

		series = append(series, s)
	}

	graph := chart.Chart{
		Width:  c.Width,
		Height: c.Height,
		DPI:    200,
		Title:  c.Title,
		XAxis: chart.XAxis{
			Name:      c.XAxisName,
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Name:      c.YAxisName,
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		Series: series,
	}

	var buffer bytes.Buffer

	err := graph.Render(chart.PNG, &buffer)

	return buffer.Bytes(), err
}
