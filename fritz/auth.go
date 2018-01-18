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

import (
    "net/url"
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/xml"
    "crypto/md5"
    "strings"
    "strconv"
    "errors"
)

// Creates a new Client instance.
// The Endpoint has to be the hostname or ip address of the Fritz!Box, HTTPS is currently not supported
// Password and Username have to be empty if there is no authentication
// Username has to be empty if a password is used for authentication
func NewClient(endpoint, username, password string) *Session {
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
func (c *Session) Login() error {
    initialSession, err := c.getInitialSessionInfo()
    if err != nil {
        return err
    }
    //Login is done at this point, because there is no authentication
    if initialSession.Valid() {
        c.sessionInfo = initialSession
        return nil
    }

    responseString := fmt.Sprintf("%s-%s", initialSession.Challenge, hashChallenge(initialSession.Challenge, c.Password))

    err = c.internalLogin(responseString)

    if err != nil {
        return err
    }

    return nil
}

func (c *Session) internalLogin(challengeResponse string) error {
    p, err := c.getLoginUrl()
    if err != nil {
        return err
    }

    data := url.Values{}
    data.Add("username", c.Username)
    data.Add("response", challengeResponse)

    req, err := http.NewRequest("POST", p.String(), strings.NewReader(data.Encode()))
    if err != nil {
        return err
    }
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

    resp, err := c.client.Do(req)
    if err != nil {
        return err
    }
    info, err := parseFromResponse(resp)
    if !info.Valid() {
        return errors.New("fritz_auth: Login Failed. Check your credentials")
    }
    c.sessionInfo = info
    return nil
}

func (c *Session) getInitialSessionInfo() (*SessionInfo, error) {
    p, err := c.getLoginUrl()
    if err != nil {
        return nil, err
    }

    resp, err := c.client.Get(p.String())
    if err != nil {
        return nil, err
    }
    return parseFromResponse(resp)
}

func (c *Session) getLoginUrl() (*url.URL, error) {
    return c.getUrl("/login_sid.lua")
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
