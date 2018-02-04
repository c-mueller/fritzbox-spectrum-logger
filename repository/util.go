package repository

import "strconv"

func (rk *RepositoryKey) GetIntegerValues() (y, m, d int) {
	year, _ := strconv.ParseInt(rk.Year, 10, 32)
	month, _ := strconv.ParseInt(rk.Month, 10, 32)
	day, _ := strconv.ParseInt(rk.Day, 10, 32)
	y, m, d = int(year), int(month), int(day)
	return
}
