// Copyright (c) 2018 Christian Müller <cmueller.dev@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package fritz

import (
    "github.com/fogleman/gg"
    "bytes"
    "image/png"
    "fmt"
    "image/color"
)

const barWidth = 4
const gridCount = 20

var gray = color.RGBA{R: 200, G: 200, B: 200, A: 0,}
var darkGray = color.RGBA{R: 20, G: 20, B: 20, A: 0,}
var black = color.RGBA{R: 0, G: 0, B: 0, A: 0,}
var purple = color.RGBA{R: 255, G: 0, B: 255, A: 0,}
var green = color.RGBA{R: 0, G: 255, B: 0, A: 0,}
var blue = color.RGBA{R: 0, G: 0, B: 255,}

func (s *Spectrum) Render() ([]byte, error) {
    w, h := s.computeSize()
    fmt.Println(w, h)
    img := gg.NewContext(w, h)
    img.SetFillRuleEvenOdd()
    img.SetLineWidth(2)

    //Fill the background white
    img.SetColor(color.White)
    img.DrawRectangle(0, 0, float64(w), float64(h))
    img.Fill()

    for index, port := range s.Ports {
        // Draw Port Number
        img.SetColor(color.Black)
        img.DrawString(fmt.Sprintf("Port #%d - %s", index, port.SpectrumInfo.ConnectionMode), 10, float64(10+index*560))

        port.renderSpectrum(index, img, port.SpectrumInfo.CurrentSNRValues,
            30, 20, renderConfig{
                PrimaryColor:   purple,
                SecondaryColor: purple,
                SecondaryAreas: make([]UpstreamRange, 0),
            })
        port.renderSpectrum(index, img, port.SpectrumInfo.CurrentBitValues,
            30, 300, renderConfig{
                PrimaryColor:   blue,
                SecondaryColor: green,
                SecondaryAreas: port.SpectrumInfo.UpstreamRanges,
            })
    }

    outputBuffer := bytes.NewBuffer([]byte(""))
    err := png.Encode(outputBuffer, img.Image())

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

    renderGrid(img, startX, startY, length)

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

func renderGrid(img *gg.Context, startX, startY, length float64) {
    for i := 1; i <= (gridCount - 1); i++ {
        setColor(img, darkGray)
        horizontalX, horizontalY := startX+float64(i)*(length/gridCount)*barWidth, startY
        img.DrawLine(horizontalX, horizontalY, horizontalX, horizontalY+260)
        img.Stroke()
        verticalX, verticalY := startX-10, startY+((float64(i)/gridCount)*250)
        img.DrawLine(verticalX, verticalY, verticalX+length*barWidth+10, verticalY)
        img.Stroke()
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
