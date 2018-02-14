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

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestDownloadSpectrum(t *testing.T) {
	spectrum, err := os.Open("testdata/example_spectrum_nt.json")
	assert.NoError(t, err)
	spectrumResponseData, err := ioutil.ReadAll(spectrum)
	assert.NoError(t, err)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log("Something was Requested")
		req, _ := ioutil.ReadAll(r.Body)
		t.Log("Request: ", string(req))
		w.Write(spectrumResponseData)
	}))
	defer testServer.Close()
	str := strings.Replace(testServer.URL, "http://", "", -1)
	session := Session{
		Endpoint: str,
		Username: "",
		Password: "",
		tokenAge: 1,
		client:   http.DefaultClient,
		sessionInfo: &SessionInfo{
			SID: "1234567890",
		},
	}

	spc, err := session.GetSpectrum()
	assert.NoError(t, err)
	assert.Equal(t, 1, spc.PortCount)
}
