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
	"crypto/md5"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Creates a new Client instance.
// The Endpoint has to be the hostname or ip address of the Fritz!Box, HTTPS is currently not supported
// Password and Username have to be empty if there is no authentication
// Username has to be empty if a password is used for authentication
func NewClient(endpoint, username, password string) *Session {
	endpoint = strings.Replace(endpoint, "http://", "", -1)
	return &Session{
		Endpoint: endpoint,
		Username: username,
		Password: password,
		client:   http.DefaultClient,
	}
}

// Attempts a Login with the Credentials set in the session.
// To Perform other operations, like downloading a spectrum, you have to call this operation
// A error is returned if the credentials are wrong or the Login could not get performed.
func (s *Session) Login() error {
	initialSession, err := s.getInitialSessionInfo()
	if err != nil {
		return err
	}
	//Login is done at this point, because there is no authentication
	if initialSession.Valid() {
		s.sessionInfo = initialSession
		s.TokenAge = time.Now().Unix()
		return nil
	}

	responseString := fmt.Sprintf("%s-%s", initialSession.Challenge, hashChallenge(initialSession.Challenge, s.Password))

	err = s.internalLogin(responseString)

	if err != nil {
		return err
	}

	return nil
}

func (s *Session) internalLogin(challengeResponse string) error {
	p, err := s.getLoginUrl()
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Add("username", s.Username)
	data.Add("response", challengeResponse)

	req, err := http.NewRequest("POST", p.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	info, err := parseFromResponse(resp)
	if !info.Valid() {
		return errors.New("fritz_auth: Login Failed. Check your credentials")
	}
	s.sessionInfo = info
	s.TokenAge = time.Now().Unix()
	return nil
}

func (s *Session) getInitialSessionInfo() (*SessionInfo, error) {
	p, err := s.getLoginUrl()
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Get(p.String())
	if err != nil {
		return nil, err
	}
	return parseFromResponse(resp)
}

func (s *Session) getLoginUrl() (*url.URL, error) {
	return s.getUrl("/login_sid.lua")
}

func parseFromResponse(resp *http.Response) (*SessionInfo, error) {
	data, err := ioutil.ReadAll(resp.Body)

	var info *SessionInfo

	err = xml.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func hashChallenge(challenge, password string) string {
	responseUnhashed := fmt.Sprintf("%s-%s", challenge, password)
	utf16EncodedResponse := convertUTF8ToUTF16LE([]byte(responseUnhashed))
	m := md5.Sum(utf16EncodedResponse)
	return fmt.Sprintf("%x", m)
}
