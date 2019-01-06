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
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	input, err := os.Open("testdata/sample_connection_info.html")
	assert.NoError(t, err)
	defer input.Close()
	data, err := ioutil.ReadAll(input)
	assert.NoError(t, err)

	conInfo, err := ParseConnectionInformation(string(data))
	assert.NoError(t, err)

	assert.Equal(t, 24984, conInfo.Downstream.CurrentDataRate)
	assert.Equal(t, 15846, conInfo.Upstream.Capacity)

}

func TestParser_WithErrs(t *testing.T) {
	input, err := os.Open("testdata/sample_connection_info_with_errs.html")
	assert.NoError(t, err)
	defer input.Close()
	data, err := ioutil.ReadAll(input)
	assert.NoError(t, err)

	conInfo, err := ParseConnectionInformation(string(data))
	assert.NoError(t, err)

	assert.Equal(t, 21112, conInfo.Downstream.CurrentDataRate)
	assert.Equal(t, 15954, conInfo.Upstream.Capacity)

}

func TestParser_EmptyInput(t *testing.T) {
	ci, err := ParseConnectionInformation("")
	assert.Error(t, err)
	assert.Nil(t, ci)
}
