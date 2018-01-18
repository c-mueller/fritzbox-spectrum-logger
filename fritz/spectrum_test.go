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
    "testing"
    "os"
    "encoding/json"
    "io/ioutil"
    "github.com/Flaque/filet"
    "fmt"
    "path/filepath"
)

func TestDrawSpectrum(t *testing.T) {
    tmpDir := filet.TmpDir(t, "")
    t.Log("Using tmpdir", tmpDir)
    //defer filet.CleanUp(t)

    t.Log("Loading test Data")
    data := loadTestData(t)
    imgdata, _ := data.Render()
    fmt.Println(len(imgdata))

    path := filepath.Join(tmpDir, "test.png")
    file, _ := os.Create(path)
    file.Write(imgdata)
    file.Close()
}

func loadTestData(t *testing.T) *Spectrum {
    file, err := os.Open("testdata/example_spectrum.json")
    if err != nil {
        t.Log(err)
        t.FailNow()
    }
    var result *Spectrum
    data, err := ioutil.ReadAll(file)
    file.Close()
    err = json.Unmarshal(data, &result)
    return result
}