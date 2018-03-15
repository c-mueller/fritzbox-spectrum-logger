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
	"github.com/fogleman/gg"
	"image/color"
)

func (p *SpectrumPort) toSNRSpectrum(portIdx int) *spectrumGraph {
	graph := spectrumGraph{
		Maximum:        p.SpectrumInfo.MaximumSNRValues,
		Minimum:        p.SpectrumInfo.MinimumSNRValues,
		Current:        p.SpectrumInfo.CurrentSNRValues,
		UpstreamRanges: make([]UpstreamRange, 0),
		PilotIndex:     -1,
		PortIndex:      portIdx,
		RenderConfig:   snrRenderConfig,
	}
	return &graph
}

func (p *SpectrumPort) toBitSpectrum(portIdx int) *spectrumGraph {
	graph := spectrumGraph{
		Maximum:        p.SpectrumInfo.MaximumBitValues,
		Minimum:        p.SpectrumInfo.MinimumBitValues,
		Current:        p.SpectrumInfo.CurrentBitValues,
		UpstreamRanges: p.SpectrumInfo.UpstreamRanges,
		PilotIndex:     p.SpectrumInfo.PilotToneIndex,
		PortIndex:      portIdx,
		RenderConfig:   bitRenderConfig,
	}
	return &graph
}

func (c ValueList) getMax() float64 {
	max := 0
	for _, v := range c {
		if v > max {
			max = v
		}
	}

	return float64(max)
}

func (s SpectrumPorts) getMaxCount() int {
	maxLen := 0
	for _, v := range s {
		if len(v.SpectrumInfo.CurrentBitValues) > maxLen {
			maxLen = len(v.SpectrumInfo.CurrentBitValues)
		}
	}

	return maxLen
}

func (g *spectrumGraph) drawPilot() bool {
	return g.PilotIndex != -1
}

func (g *spectrumGraph) useSecondary(idx int) bool {
	for _, v := range g.UpstreamRanges {
		if idx >= v.FirstIndex && idx <= v.LastIndex {
			return true
		}
	}
	return false
}


func setColor(img *gg.Context, color color.RGBA) {
	img.SetRGB255(int(color.R), int(color.G), int(color.B))
}

func (s *Spectrum) computeSize() (int, int) {
	height := s.PortCount*(30+2*maxSpectrumHeight+30) + 50
	width := 60 + barWidth*s.Ports.getMaxCount()
	return width, height
}