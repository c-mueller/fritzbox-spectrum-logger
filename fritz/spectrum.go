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
    "io/ioutil"
    "fmt"
    "encoding/json"
    "errors"
    "time"
)

var errSessionTimedOut = errors.New("spectrum_dl: Downloading spectrum failed, maybe the session timed out")

// This operation downloads the Spectrum from the Fritz!Box
// A error gets returned if the download fails.
func (c *Session) GetSpectrum() (*Spectrum, error) {
    spcUrl, err := c.getSpectrumUrl()
    if err != nil {
        return nil, err
    }

    res, err := c.client.Get(spcUrl.String())
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

    return spectrum, nil
}

func (c *Session) getSpectrumUrl() (*url.URL, error) {
    return c.getUrl(fmt.Sprintf("/internet/dsl_spectrum.lua?sid=%s&useajax=1", c.sessionInfo.SID))
}

// Converts the given spectrum to a ByteArray
// Returns a error if the Marshalling fails
func (s *Spectrum) JSON() ([]byte, error) {
    return json.Marshal(s)
}
