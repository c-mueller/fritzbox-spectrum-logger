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
	"bytes"
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/fogleman/gg"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image/color"
	"image/png"
	"math"
	"time"
)

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
	fontBox, err := rice.FindBox(fontBoxName)
	if err != nil {
		return nil, err
	}

	fontBytes, err := fontBox.Bytes(fontPath)
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

		//TODO Support Multiple Ports

		//Render SNR Spectrum
		port.toSNRSpectrum(index).render(30, 20, img)

		//Render Bit Spectrum
		port.toBitSpectrum(index).render(30, 300, img)
	}

	outputBuffer := bytes.NewBuffer([]byte(""))
	err = png.Encode(outputBuffer, img.Image())

	if err != nil {
		return nil, err
	}

	return outputBuffer.Bytes(), nil
}

func (g *spectrumGraph) render(startX, startY float64, img *gg.Context) {
	setColor(img, gray)
	img.DrawRectangle(startX, startY+(float64(g.PortIndex)*560), float64(len(g.Current)*barWidth), maxSpectrumHeight)
	img.Fill()

	maxHeight := float64(g.Current.getMax() * 1.10)

	g.renderGrid(startX, startY, img)

	for idx, valueHeight := range g.Current {
		if g.useSecondary(idx) {
			setColor(img, g.RenderConfig.SecondaryColor)
		} else {
			setColor(img, g.RenderConfig.PrimaryColor)
		}
		height := (float64(valueHeight) / maxHeight) * float64(maxSpectrumHeight)
		x, y := startX+float64(idx)*barWidth, startY+(float64(g.PortIndex*560))+(maxSpectrumHeight-height)
		img.DrawRectangle(x, y, barWidth, height)
		img.Fill()
	}

	if g.drawPilot() {
		setColor(img, g.RenderConfig.PilotColor)
		x, y := startX+float64(g.PilotIndex)*barWidth, startY+(float64(g.PortIndex*560))
		img.DrawRectangle(x, y, barWidth, maxSpectrumHeight)
		img.Fill()
	}
}

func (g *spectrumGraph) renderGrid(startX, startY float64, img *gg.Context) {
	//Calculate Maximum Horizontal Value
	length := float64(len(g.Current))
	max := g.Current.getMax() * 1.10
	scale := max / float64(horizontalGridCount)
	if scale == 0 {
		scale = 1
	}

	setColor(img, darkGray)
	//Draw Horizontal Lines of the grid
	for i := 1; i <= (horizontalGridCount - 1); i++ {
		hX, hY := startX-10, startY+((float64(i)/horizontalGridCount)*maxSpectrumHeight)
		img.DrawLine(hX, hY, hX+length*barWidth+gridLineOffset, hY)
		img.Stroke()
		hLineText := int(math.Ceil(scale * float64(horizontalGridCount-i)))
		img.DrawString(fmt.Sprintf("%d", hLineText), hX-20, hY)
	}

	//Draw Vertical lines of the grid
	for i := 1; i <= (verticalGridCount - 1); i++ {
		vX, vY := startX+float64(i)*(length/verticalGridCount)*barWidth, startY
		img.DrawLine(vX, vY, vX, vY+maxSpectrumHeight+gridLineOffset)
		img.Stroke()
		//Draw Line Description
		vLineText := int64(math.Ceil(length * (float64(i) / verticalGridCount)))
		img.DrawString(fmt.Sprintf("%d", vLineText), vX+5, vY+maxSpectrumHeight+gridLineOffset+5)
	}
}

func setColor(img *gg.Context, color color.RGBA) {
	img.SetRGB255(int(color.R), int(color.G), int(color.B))
}

func (s *Spectrum) computeSize() (int, int) {
	height := s.PortCount*(30+2*maxSpectrumHeight+30) + 50
	width := 60 + barWidth*s.Ports.getMaxCount()
	return width, height
}
