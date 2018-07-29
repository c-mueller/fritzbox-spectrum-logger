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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"
)

var errSessionTimedOut = errors.New("spectrum_dl: Downloading spectrum failed, maybe the session timed out")

// This operation downloads the Spectrum from the Fritz!Box
// A error gets returned if the download fails.
func (s *Session) GetSpectrum() (*Spectrum, error) {
	specUrl, err := s.getSpectrumUrl()
	if err != nil {
		return nil, err
	}

	res, err := s.client.Get(specUrl.String())
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(res.Body)

	var spectrum *Spectrum

	err = json.Unmarshal(data, &spectrum)
	if err != nil {
		return nil, errSessionTimedOut
	} else if spectrum.PortCount < 1 {
		return nil, errSessionTimedOut
	}

	spectrum.Timestamp = time.Now().Unix()

	conInfoUrl, err := s.getConnectionInfoUrl()
	if err != nil {
		return nil, err
	}
	res, err = s.client.Get(conInfoUrl.String())

	conInfoData, _ := ioutil.ReadAll(res.Body)

	spectrum.ConnectionInformation = string(conInfoData)

	return spectrum, nil
}

func (s *Session) getSpectrumUrl() (*url.URL, error) {
	return s.getUrl(fmt.Sprintf("/internet/dsl_spectrum.lua?sid=%s&useajax=1", s.sessionInfo.SID))
}

func (s *Session) getConnectionInfoUrl() (*url.URL, error) {
	return s.getUrl(fmt.Sprintf("/internet/dsl_stats_tab.lua?update=mainDiv&sid=%s", s.sessionInfo.SID))
}

// Converts the given spectrum to a ByteArray
// Returns a error if the Marshalling fails
func (s *Spectrum) JSON() ([]byte, error) {
	return json.Marshal(s)
}
