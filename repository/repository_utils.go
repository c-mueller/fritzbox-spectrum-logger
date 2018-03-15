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
	"github.com/boltdb/bolt"
	"github.com/c-mueller/fritzbox-spectrum-logger/fritz"
)

func (r *Repository) GetSpectraForSpectrumKey(k SpectrumKey) ([]*fritz.Spectrum, error) {
	if !k.IsValid() {
		return nil, InvalidDateKey
	}
	y, m, d := k.GetIntegerValues()
	return r.GetSpectraForDay(d, m, y)
}

func (r *Repository) GetSpectrumBySpectrumKey(k *SpectrumKey, timestamp int64) (*fritz.Spectrum, error) {
	if !k.IsValid() {
		return nil, InvalidDateKey
	}
	y, m, d := k.GetIntegerValues()
	return r.GetSpectrum(d, m, y, timestamp)
}

func (r *Repository) forEachSpectrumInDay(y, m, d int, operator func(dayBucket *bolt.Bucket, k, v []byte) error) error {
	yearByte, monthByte, dayByte := convertToByte(y, m, d)

	err := r.db.View(func(tx *bolt.Tx) error {

		dayBucket, err := r.getDayBucket(dayByte, monthByte, yearByte, tx)
		if err != nil {
			return err
		}

		err = dayBucket.ForEach(func(k, v []byte) error {
			return operator(dayBucket, k, v)
		})
		return err
	})

	return err
}

func (r *Repository) forEachSpectrumKey(operator func(dayBucket *bolt.Bucket, key SpectrumKey) error) error {
	err := r.db.View(func(tx *bolt.Tx) error {
		spectraBucket := tx.Bucket([]byte(SpectrumListBucketName))

		err := spectraBucket.ForEach(func(yearKey, v []byte) error {
			yearBucket := spectraBucket.Bucket(yearKey)
			if yearBucket == nil {
				//Ignore element if it is not a bucket
				return nil
			}
			err := yearBucket.ForEach(func(monthKey, v []byte) error {
				monthBucket := yearBucket.Bucket(monthKey)
				if monthBucket == nil {
					//Ignore element if it is not a bucket
					return nil
				}
				err := monthBucket.ForEach(func(dayKey, v []byte) error {
					dayBucket := monthBucket.Bucket(dayKey)
					if dayBucket == nil {
						//Ignore element if it is not a bucket
						return nil
					}
					key := SpectrumKey{
						Year:  string(yearKey),
						Month: string(monthKey),
						Day:   string(dayKey),
					}
					return operator(dayBucket, key)
				})
				return err
			})
			return err
		})
		return err
	})
	return err
}

func (r *Repository) getDayBucket(dayByte, monthByte, yearByte []byte, tx *bolt.Tx) (*bolt.Bucket, error) {
	spectraBucket := tx.Bucket([]byte(SpectrumListBucketName))
	if spectraBucket == nil {
		log.Error("Spectra Bucket not found!")
		return nil, BucketNotFoundError
	}
	yearBucket := spectraBucket.Bucket(yearByte)
	if yearBucket == nil {
		log.Errorf("Year Bucket (Year: '%s') not found!",
			string(yearByte))
		return nil, BucketNotFoundError
	}
	monthBucket := yearBucket.Bucket(monthByte)
	if monthBucket == nil {
		log.Errorf("Month Bucket (Year: '%s' Month: '%s') not found!",
			string(yearByte), string(monthByte))
		return nil, BucketNotFoundError
	}
	dayBucket := monthBucket.Bucket(dayByte)
	if dayBucket == nil {
		log.Errorf("Month Bucket (Year: '%s' Month: '%s' Day: '%s') not found!",
			string(yearByte), string(monthByte), string(dayByte))
		return nil, BucketNotFoundError
	}
	return dayBucket, nil
}
