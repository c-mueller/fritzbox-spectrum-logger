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
	"unicode/utf16"
	"unicode/utf8"
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
