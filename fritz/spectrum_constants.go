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

import "image/color"

const barWidth = 4
const maxSpectrumHeight = 250
const gridLineOffset = 10
const horizontalGridCount = 4
const verticalGridCount = 8

const fontBoxName = "font"
const fontPath = "luxisr.ttf"

var gray = color.RGBA{R: 200, G: 200, B: 200, A: 0}
var darkGray = color.RGBA{R: 20, G: 20, B: 20, A: 0}
var black = color.RGBA{R: 0, G: 0, B: 0, A: 0}
var red = color.RGBA{R: 255, G: 0, B: 0, A: 0}
var purple = color.RGBA{R: 255, G: 0, B: 255, A: 0}
var green = color.RGBA{R: 0, G: 255, B: 0, A: 0}
var blue = color.RGBA{R: 0, G: 0, B: 255}

var bitRenderConfig = renderConfig{
	PrimaryColor:   blue,
	SecondaryColor: green,
	PilotColor:     purple,
	MinColor:       red,
	MaxColor:       red,
}

var snrRenderConfig = renderConfig{
	PrimaryColor:   purple,
	SecondaryColor: purple,
	MinColor:       red,
	MaxColor:       red,
}
