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
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestFritzAuthNoPasswd(t *testing.T) {
	validFile, err := os.Open("testdata/auth_valid.xml")
	assert.NoError(t, err)
	response, err := ioutil.ReadAll(validFile)
	assert.NoError(t, err)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log("Something was Requested")
		req, _ := ioutil.ReadAll(r.Body)
		t.Log("Request: ", string(req))
		w.Write(response)
	}))
	defer testServer.Close()

	session := NewClient(testServer.URL, "", "")
	err = session.Login()
	assert.NoError(t, err)
	assert.True(t, session.Valid())
}

func TestFritzAuthPassword_Fail(t *testing.T) {
	invalidFile, err := os.Open("testdata/auth_invalid.xml")
	assert.NoError(t, err)
	response, err := ioutil.ReadAll(invalidFile)
	assert.NoError(t, err)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log("Something was Requested")
		req, _ := ioutil.ReadAll(r.Body)
		t.Log("Request: ", string(req))
		w.Write(response)
	}))
	defer testServer.Close()

	session := NewClient(testServer.URL, "", "123456")
	err = session.Login()
	assert.Error(t, err)
	assert.True(t, !session.Valid())
}

func TestFritzAuthPassword_Success(t *testing.T) {
	invalidFile, err := os.Open("testdata/auth_invalid.xml")
	assert.NoError(t, err)
	invalidResponse, err := ioutil.ReadAll(invalidFile)
	assert.NoError(t, err)

	validFile, err := os.Open("testdata/auth_valid.xml")
	assert.NoError(t, err)
	validResponse, err := ioutil.ReadAll(validFile)
	assert.NoError(t, err)

	calledOnce := false

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log("Something was Requested")
		req, _ := ioutil.ReadAll(r.Body)
		t.Log("Request: ", string(req))
		if !calledOnce {
			t.Log("Responded With Inalid Response")
			w.Write(invalidResponse)
			calledOnce = true
		} else {
			t.Log("Responded With Valid Response")
			w.Write(validResponse)
		}
	}))
	defer testServer.Close()

	session := NewClient(testServer.URL, "", "123456")
	err = session.Login()
	assert.NoError(t, err)
	assert.True(t, session.Valid())
}
