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
    "unicode/utf8"
    "unicode/utf16"
)

func convertUTF8ToUTF16LE(p []byte) []byte {
    bs := make([]byte, 0, 2*len(p))
    pos := 0
    for pos < len(p) {
        bytes, size := consumeNextRune(p[pos:])
        pos += size
        bs = append(bs, bytes...)
    }
    return bs
}

func consumeNextRune(p []byte) ([]byte, int) {
    r, size := utf8.DecodeRune(p)
    if r <= 0xffff {
        return []byte{uint8(r), uint8(r >> 8)}, size
    }
    r1, r2 := utf16.EncodeRune(r)
    return []byte{uint8(r1), uint8(r1 >> 8), uint8(r2), uint8(r2 >> 8)}, size
}
