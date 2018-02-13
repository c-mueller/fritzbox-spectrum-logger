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

type SpectrumPorts []SpectrumPort
type ValueList []int

type Spectrum struct {
	PortCount int           `json:"ports"`
	Ports     SpectrumPorts `json:"port"`
	Timestamp int64         `json:"timestamp"`
}

type SpectrumPort struct {
	SpectrumInfo SpectrumInfo `json:"us"`
}

type SpectrumInfo struct {
	TonesPerBATValue    int             `json:"TONES_PER_BAT_VALUE"`
	MaximumSNRFrequency int             `json:"MAX_SNR_FREQ"`
	PilotToneIndex      int             `json:"PILOT"`
	UpstreamRanges      []UpstreamRange `json:"BIT_BANDCONFIG"`
	DetectedNoiseValues ValueList       `json:"DETECTED_RFI_VALUES"`
	ConnectionMode      string          `json:"MODE"`
	MaximumBATFrequency int             `json:"MAX_BAT_FREQ"`
	TonesPerSNRValue    int             `json:"TONES_PER_SNR_VALUE"`
	CurrentBitValues    ValueList       `json:"ACT_BIT_VALUES"`
	MaximumBitValues    ValueList       `json:"MAX_BIT_VALUES"`
	MinimumBitValues    ValueList       `json:"MIN_BIT_VALUES"`
	CurrentSNRValues    ValueList       `json:"ACT_SNR_VALUES"`
	MaximumSNRValues    ValueList       `json:"MAX_SNR_VALUES"`
	MinimumSNRValues    ValueList       `json:"MIN_SNR_VALUES"`
}

type UpstreamRange struct {
	FirstIndex int `json:"FIRST"`
	LastIndex  int `json:"LAST"`
}

type renderConfig struct {
	PrimaryColor   color.RGBA
	SecondaryColor color.RGBA
	SecondaryAreas []UpstreamRange
}
