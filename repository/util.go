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

package repository

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

func GetFromTimestamp(timestamp int64) SpectrumKey {
	t := time.Unix(timestamp, 0)
	return SpectrumKey{
		Day:   fmt.Sprintf("%d", t.Day()),
		Month: fmt.Sprintf("%d", int(t.Month())),
		Year:  fmt.Sprintf("%d", t.Year()),
	}
}

func (sk *SpectrumKey) GetIntegerValues() (y, m, d int) {
	year, err := strconv.ParseInt(sk.Year, 10, 32)
	if err != nil {
		return -1, -1, -1
	}
	month, err := strconv.ParseInt(sk.Month, 10, 32)
	if err != nil {
		return -1, -1, -1
	}
	day, err := strconv.ParseInt(sk.Day, 10, 32)
	if err != nil {
		return -1, -1, -1
	}
	y, m, d = int(year), int(month), int(day)
	return
}

func (sk *SpectrumKey) IsValid() bool {
	y, _, _ := sk.GetIntegerValues()
	return y != -1
}

func (sk *SpectrumKey) String() string {
	y, m, d := sk.GetIntegerValues()
	return fmt.Sprintf("%02d.%02d.%04d", d, m, y)
}

func (ta TimestampArray) Search(value int64) int {
	return sort.Search(ta.Len(), func(i int) bool {
		return ta[i] >= value
	})
}

func (ta TimestampArray) Len() int {
	return len(ta)
}

func (ta TimestampArray) Swap(i, j int) {
	ta[i], ta[j] = ta[j], ta[i]
}

func (ta TimestampArray) Less(i, j int) bool {
	return ta[i] < ta[j]
}

func (k SpectraKeys) Len() int {
	return len(k)
}

func (k SpectraKeys) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}

func (k SpectraKeys) Less(i, j int) bool {
	aY, aM, aD := k[i].GetIntegerValues()
	bY, bM, bD := k[j].GetIntegerValues()
	if aY == bY {
		if aM == bM {
			if aD == bD {
				return false
			} else {
				return aD < bD
			}
		} else {
			return aM < bM
		}
	} else {
		return aY < bY
	}
}
