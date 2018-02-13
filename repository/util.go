package repository

import (
	"strconv"
)

func (sk *SpectrumKey) GetIntegerValues() (y, m, d int) {
	year, _ := strconv.ParseInt(sk.Year, 10, 32)
	month, _ := strconv.ParseInt(sk.Month, 10, 32)
	day, _ := strconv.ParseInt(sk.Day, 10, 32)
	y, m, d = int(year), int(month), int(day)
	return
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
