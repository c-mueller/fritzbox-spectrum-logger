// Copyright (c) 2018 Christian Müller <cmueller.dev@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
	return s.SID != "0000000000000000"
}
