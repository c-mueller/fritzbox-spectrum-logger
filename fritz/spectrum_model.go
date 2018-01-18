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
