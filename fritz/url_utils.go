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
	"fmt"
	"net/url"
)

func (c *Session) getUrl(rawurl string) (*url.URL, error) {
	base, err := url.Parse(fmt.Sprintf("http://%s", c.Endpoint))
	if err != nil {
		c.client = nil
		return nil, err
	}
	loginPath, err := url.Parse(rawurl)
	if err != nil {
		c.client = nil
		return nil, err
	}

	return base.ResolveReference(loginPath), nil
}
