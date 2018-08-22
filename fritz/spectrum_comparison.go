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
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"sort"
	"time"
)

const comparisonPortIndex = 0

func (c ComparisonSet) RenderComparison(scaled bool) ([]byte, error) {
	sort.Sort(c)

	w, h := c.computeComparisonDimensions(scaled)
	img, err := initializeImageContext(w, h, scaled)
	if err != nil {
		return nil, err
	}

	maxBit := c.getBitMaxHeight()
	maxSNR := c.getSNRMaxHeight()
	graphWidth := c.getMaxEntryCount()

	for i, v := range c {
		port := v.Ports[0]
		height := 20 + float64(i*(30+maxSpectrumHeight+30))

		img.SetColor(color.Black)
		spectrumTime := time.Unix(v.Timestamp, 0).String()
		portHeading := fmt.Sprintf("Measurement #%d from %s", i+1, spectrumTime)

		tW, tH := img.MeasureString(portHeading)

		img.DrawString(portHeading, (float64(w)-tW)/2, height-(tH/2))

		port.toSNRSpectrum(0).render(30, height, maxSNR, img)
		port.toBitSpectrum(0).render(30+float64(graphWidth*barWidth+70), height, maxBit, img)
	}

	outputBuffer := bytes.NewBuffer([]byte(""))
	err = png.Encode(outputBuffer, img.Image())

	if err != nil {
		return nil, err
	}

	return outputBuffer.Bytes(), nil
}

func (c ComparisonSet) Len() int {
	return len(c)
}
func (c ComparisonSet) Swap(a, b int) {
	c[a], c[b] = c[b], c[a]
}

func (c ComparisonSet) Less(a, b int) bool {
	return c[a].Timestamp < c[b].Timestamp
}
