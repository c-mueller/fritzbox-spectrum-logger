package repository

import (
	"sort"
	"time"
)

const maxFails = 365

func GetNeighbours(r Repository, timestamp int64) (int64, int64, error) {
	timestamps, err := r.GetTimestampsForSpectrumKey(GetFromTimestamp(timestamp))
	if err != nil {
		return -1, -1, err
	}

	sort.Sort(timestamps)

	pos := sort.Search(len(timestamps), func(i int) bool {
		return timestamps[i] >= timestamp
	})

	if pos != 0 && pos != (len(timestamps)-1) {
		return timestamps[pos-1], timestamps[pos+1], nil
	} else if pos == 0 {
		prevLen := 0
		iteration := maxFails
		selectedTs := int64(-1)
		for prevLen == 0 && iteration > 0 {
			day := time.Unix(timestamp, 0).Add(-1 * time.Duration(maxFails-iteration+1) * time.Hour * 24)
			ts, _ := r.GetTimestampsForSpectrumKey(GetFromTimestamp(day.Unix()))

			prevLen = len(ts)

			if prevLen != 0 {
				selectedTs = ts[len(ts)-1]
			}

			iteration--
		}

		if pos+1 >= len(timestamps) {
			return selectedTs, -1, nil
		}

		return selectedTs, timestamps[pos+1], nil
	} else if pos == len(timestamps)-1 {
		prevLen := 0
		iteration := maxFails
		selectedTs := int64(-1)
		for prevLen == 0 && iteration > 0 {
			day := time.Unix(timestamp, 0).Add(time.Duration(maxFails-iteration+1) * time.Hour * 24)
			ts, _ := r.GetTimestampsForSpectrumKey(GetFromTimestamp(day.Unix()))

			prevLen = len(ts)

			if prevLen != 0 {
				selectedTs = ts[0]
			}

			iteration--
		}

		if pos-1 < 0 {
			return -1, selectedTs, nil
		}

		return timestamps[pos-1], selectedTs, nil
	}

	return -1, -1, nil
}
