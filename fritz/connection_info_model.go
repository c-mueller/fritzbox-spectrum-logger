// Fritz!Box Spectrum Logger (https://github.com/c-mueller/fritzbox-spectrum-logger).
// Copyright (c) 2018 Christian Müller <cmueller.dev@gmail.com>.
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

type ConnectionInformation struct {
	Downstream ConnectionTransmissionDirection `json:"downstream"`
	Upstream   ConnectionTransmissionDirection `json:"upstream"`
	Profile    string                          `json:"profile"`
}

type ConnectionTransmissionDirection struct {
	MinimumDataRate        int  `json:"minimum_data_rate"`
	MaximumDataRate        int  `json:"maximum_data_rate"`
	Capacity               int  `json:"capacity"`
	CurrentDataRate        int  `json:"current_data_rate"`
	SeamlessRateAdjustment bool `json:"seamless_rate_adjustment"`

	Latency  int     `json:"latency"`
	INPValue float64 `json:"inp_value"`
	GINP     bool    `json:"ginp"`

	SNMargin        float64 `json:"sn_margin"`
	Bitswap         bool    `json:"bitswap"`
	LineAttenuation float64 `json:"line_attenuation"`

	VectorMode string `json:"vector_mode"`
	Carrier    string `json:"carrier"`

	Errors Errors `json:"errors"`
}

type Errors struct {
	SecondsWithErrors     float64 `json:"seconds_with_errors"`
	SecondsWithManyErrors float64 `json:"seconds_with_many_errors"`
	ErrorsPerMinute       float64 `json:"errors_per_minute"`
	ErrorsLast15Min       float64 `json:"errors_last_15_min"`
}
