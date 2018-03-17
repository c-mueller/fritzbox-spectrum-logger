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
	"github.com/fogleman/gg"
	"image/color"
	"image/png"
	"math"
	"time"
	"github.com/GeertJohan/go.rice"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

// Renders the spectrum to a PNG byte array
func (s *Spectrum) Render(scaled bool) ([]byte, error) {
	w, h := s.computeSize()
	img := gg.NewContext(w, h)
	if scaled {
		img = gg.NewContext(w*2, h*2)
		img.Scale(2, 2)
	}

	img.SetFillRuleEvenOdd()
	img.SetLineCapSquare()
	img.SetLineWidth(1.0)

	//Fill the background white
	img.SetColor(color.White)
	img.DrawRectangle(0, 0, float64(w), float64(h))
	img.Fill()

	if scaled {
		//Set Font
		fontBox, err := rice.FindBox("font")
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
	}

	for index, port := range s.Ports {
		// Draw Port Number
		img.SetColor(color.Black)
		spectrumTime := time.Unix(s.Timestamp, 0).String()
		portHeading := fmt.Sprintf("Port #%d - %s - (From %s)", index, port.SpectrumInfo.ConnectionMode, spectrumTime)
		img.DrawString(portHeading, 10, float64(15+index*560))

		//Render SNR Spectrum
		port.toSNRSpectrum(index).render(30, 20+float64(index*(30+2*maxSpectrumHeight+30)), img)
		//Render Bit Spectrum
		port.toBitSpectrum(index).render(30, 300+float64(index*(30+2*maxSpectrumHeight+30)), img)
	}

	outputBuffer := bytes.NewBuffer([]byte(""))
	err := png.Encode(outputBuffer, img.Image())

	if err != nil {
		return nil, err
	}

	return outputBuffer.Bytes(), nil
}

func (g *spectrumGraph) render(startX, startY float64, img *gg.Context) {
	g.fillBackground(startX, startY, img)

	g.renderGrid(startX, startY, img)

	g.renderCurrent(startX, startY, img)

	//Draw Minimum and maxiumum
	g.renderLine(g.Minimum, g.RenderConfig.MinColor, startX, startY, img) // Minimum
	g.renderLine(g.Maximum, g.RenderConfig.MaxColor, startX, startY, img) // Maximum

	g.renderPilot(startX, startY, img)
}

func (g *spectrumGraph) renderLine(data ValueList, lineColor color.RGBA, startX, startY float64, img *gg.Context) {
	maxHeight := float64(g.Current.getMax() * 1.10)
	oldY := float64(0)
	for idx, heightValue := range data {
		setColor(img, lineColor)
		height := (float64(heightValue) / maxHeight) * float64(maxSpectrumHeight)
		x, y := startX+float64(idx)*barWidth, startY+(maxSpectrumHeight-height)
		if idx != 0 && idx-1 != g.PilotIndex {
			img.DrawRectangle(x, oldY, 1, y-oldY)
		}
		img.DrawRectangle(x, y, barWidth, 1)
		img.Fill()
		oldY = y
	}
}

func (g *spectrumGraph) renderPilot(startX, startY float64, img *gg.Context) {
	if g.drawPilot() {
		setColor(img, g.RenderConfig.PilotColor)
		x, y := startX+float64(g.PilotIndex)*barWidth, startY
		img.DrawRectangle(x, y, barWidth, maxSpectrumHeight)
		img.Fill()
	}
}

func (g *spectrumGraph) renderCurrent(startX, startY float64, img *gg.Context) {
	maxHeight := float64(g.Current.getMax() * 1.10)
	for idx, heightValue := range g.Current {
		if g.useSecondary(idx) {
			setColor(img, g.RenderConfig.SecondaryColor)
		} else {
			setColor(img, g.RenderConfig.PrimaryColor)
		}
		height := (float64(heightValue) / maxHeight) * float64(maxSpectrumHeight)
		x, y := startX+float64(idx)*barWidth, startY+(maxSpectrumHeight-height)
		if heightValue != 0 {
			img.DrawRectangle(x, y, barWidth, height)
			img.Fill()
		}
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

	setColor(img, g.RenderConfig.GridColor)
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

func (g *spectrumGraph) fillBackground(startX, startY float64, img *gg.Context) {
	setColor(img, g.RenderConfig.BackgroundColor)
	img.DrawRectangle(startX, startY, float64(len(g.Current)*barWidth), maxSpectrumHeight)
	img.Fill()
}
