// Fritz!Box Spectrum Logger (https://github.com/c-mueller/fritzbox-spectrum-logger).
// Copyright (c) 2018 Christian Müller<cmueller.dev@gmail.com>.
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
	"fmt"
	"github.com/fogleman/gg"
	"image/color"
	"image/png"
	"time"
	"github.com/GeertJohan/go.rice"
	"github.com/golang/freetype/truetype"
	"bytes"
	"github.com/golang/freetype"
)

const barWidth = 4
const gridCount = 20

var gray = color.RGBA{R: 200, G: 200, B: 200, A: 0}
var darkGray = color.RGBA{R: 20, G: 20, B: 20, A: 0}
var black = color.RGBA{R: 0, G: 0, B: 0, A: 0}
var purple = color.RGBA{R: 255, G: 0, B: 255, A: 0}
var green = color.RGBA{R: 0, G: 255, B: 0, A: 0}
var blue = color.RGBA{R: 0, G: 0, B: 255}

// Renders the spectrum to a PNG byte array
func (s *Spectrum) Render() ([]byte, error) {
	w, h := s.computeSize()

	img := gg.NewContext(w*2, h*2)
	img.Scale(2, 2)

	//img := gg.NewContext(w, h)

	img.SetFillRuleEvenOdd()
	img.SetLineCapSquare()
	img.SetLineWidth(1.0)

	//Fill the background white
	img.SetColor(color.White)
	img.DrawRectangle(0, 0, float64(w), float64(h))
	img.Fill()

	//Set Font
	fontBox, err := rice.FindBox("font")
	if err != nil {
		return nil, err
	}
	fontBytes, err := fontBox.Bytes("luxisr.ttf")
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	img.SetFontFace(truetype.NewFace(font, &truetype.Options{
		Size: 12,
		DPI:  144,
	}))

	for index, port := range s.Ports {
		// Draw Port Number
		img.SetColor(color.Black)
		spectrumTime := time.Unix(s.Timestamp, 0).String()
		portHeading := fmt.Sprintf("Port #%d - %s - (From %s)", index, port.SpectrumInfo.ConnectionMode, spectrumTime)
		img.DrawString(portHeading, 10, float64(15+index*560))

		//Render SNR Spectrum
		port.renderSpectrum(index, img, port.SpectrumInfo.CurrentSNRValues,
			30, 20, renderConfig{
				PrimaryColor:   purple,
				SecondaryColor: purple,
				SecondaryAreas: make([]UpstreamRange, 0),
			})

		//Render Bit Spectrum
		port.renderSpectrum(index, img, port.SpectrumInfo.CurrentBitValues,
			30, 300, renderConfig{
				PrimaryColor:   blue,
				SecondaryColor: green,
				SecondaryAreas: port.SpectrumInfo.UpstreamRanges,
			})
	}

	outputBuffer := bytes.NewBuffer([]byte(""))
	err = png.Encode(outputBuffer, img.Image())

	if err != nil {
		return nil, err
	}

	return outputBuffer.Bytes(), nil
}

func (port *SpectrumPort) renderSpectrum(index int, img *gg.Context, data ValueList, startX, startY float64, config renderConfig) {
	setColor(img, gray)
	img.DrawRectangle(startX, startY+(float64(index)*560), float64(len(port.SpectrumInfo.CurrentBitValues)*barWidth), 250)
	img.Fill()

	maxHeight := float64(data.getMax() * 1.10)
	length := float64(len(data))

	port.renderGrid(img, startX, startY, length, data)

	for idx, valueHeight := range data {
		if config.useSecondary(idx) {
			setColor(img, config.SecondaryColor)
		} else {
			setColor(img, config.PrimaryColor)
		}
		height := (float64(valueHeight) / maxHeight) * 250.0
		x, y := startX+float64(idx)*barWidth, startY+(float64(index * 560))+(250-height)
		img.DrawRectangle(x, y, barWidth, height)
		img.Fill()
	}

}

func (port *SpectrumPort) renderGrid(img *gg.Context, startX, startY, length float64, data ValueList) {
	//Calculate Maximum Horizontal Value
	max := data.getMax() * 1.10
	scale := int(max / gridCount)
	if scale == 0 {
		scale = 1
	}

	for i := 1; i <= (gridCount - 1); i++ {
		setColor(img, darkGray)
		//Draw Vertical lines of the grid
		vX, vY := startX+float64(i)*(length/gridCount)*barWidth, startY
		img.DrawLine(vX, vY, vX, vY+260)
		img.Stroke()
		//Draw Line Description
		vLineValue := int64(length * (float64(i) / gridCount))
		img.DrawString(fmt.Sprintf("%d", vLineValue), vX+5, vY+260+5)

		//Draw Horizontal Lines of the grid
		hX, hY := startX-10, startY+((float64(i)/gridCount)*250)
		img.DrawLine(hX, hY, hX+length*barWidth+10, hY)
		img.Stroke()
		img.DrawString(fmt.Sprintf("%d", scale*(gridCount-i)), hX-20, hY)
	}
}

func setColor(img *gg.Context, color color.RGBA) {
	img.SetRGB255(int(color.R), int(color.G), int(color.B))
}

func (s *Spectrum) computeSize() (int, int) {
	height := s.PortCount*(30+2*250+30) + 50
	width := 60 + barWidth*s.Ports.getMaxCount()
	return width, height
}
