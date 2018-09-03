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

package migrator

import (
	"github.com/c-mueller/fritzbox-spectrum-logger/repository"
	"gopkg.in/cheggaaa/pb.v1"
)

func NewMigrator(target repository.Repository, verbose bool, source ... repository.Repository) (*DatabaseMigrator, error) {
	return &DatabaseMigrator{
		TargetRepository:   target,
		SourceRepositories: source,
		Verbose:            verbose,
	}, nil
}

func (m *DatabaseMigrator) Migrate() error {
	for k, source := range m.SourceRepositories {
		log.Infof("Migrating database #%d...", k+1)
		sKeys, err := source.GetAllSpectrumKeys()
		if err != nil {
			return err
		}

		log.Infof("Found %d Days (SpectrumKeys) in database #%d...", len(sKeys), k)

		for dayIdx, spectrumKey := range sKeys {
			timestamps, err := source.GetTimestampsForSpectrumKey(spectrumKey)
			if err != nil {
				return err
			}

			log.Infof("Inserting %d Spectra for %s...", len(timestamps), spectrumKey.String())

			bar := pb.New(len(timestamps))
			bar.Start()

			bar.NotPrint = !m.Verbose

			for _, timestamp := range timestamps {
				year, month, day := spectrumKey.GetIntegerValues()
				spectrum, err := source.GetSpectrum(day, month, year, timestamp)
				if err != nil {
					return err
				}

				err = m.TargetRepository.Insert(spectrum)
				if err != nil {
					return err
				}
				bar.Increment()
			}
			bar.Finish()
			log.Infof("Completed Day %d/%d...", dayIdx, len(sKeys))
		}
	}
	return nil
}
