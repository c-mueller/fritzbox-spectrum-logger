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

import "net/http"

type Session struct {
	Endpoint    string
	Username    string
	Password    string
	sessionInfo *SessionInfo
	tokenAge    int64
	client      *http.Client
}

type SessionInfo struct {
	Challenge string     `xml:"Challenge"`
	SID       string     `xml:"SID"`
	BlockTime string     `xml:"BlockTime"`
	Rights    Privileges `xml:"Rights"`
}

type Privileges struct {
	Names        []string `xml:"Name"`
	AccessLevels []string `xml:"Access"`
}

func (s *SessionInfo) Valid() bool {
	if s == nil {
		return false
	}
	return  s.SID != "0000000000000000"
}

func (s *Session) Valid() bool {
	return s.sessionInfo.Valid()
}
