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

func (r *renderConfig) useSecondary(idx int) bool {
	for _, v := range r.SecondaryAreas {
		if idx >= v.FirstIndex && idx <= v.LastIndex {
			return true
		}
	}
	return false
}
