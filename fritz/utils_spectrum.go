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

func (c ValueList) getMax() float64 {
    max := 0
    for _, v := range c {
        if v > max {
            max = v
        }
    }

    return float64(max)
}

func (s SpectrumPorts) getMaxCount() int {
    maxLen := 0
    for _, v := range s {
        if len(v.SpectrumInfo.CurrentBitValues) > maxLen {
            maxLen = len(v.SpectrumInfo.CurrentBitValues)
        }
    }

    return maxLen
}

func (r *renderConfig) useSecondary(idx int) bool {
    for _, v := range r.SecondaryAreas{
        if idx >= v.FirstIndex && idx <= v.LastIndex {
            return true
        }
    }
    return false
}