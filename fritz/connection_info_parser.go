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

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var ParsingError = errors.New("Parsing of Connection info has failed")

var tdRegex *regexp.Regexp
var unitRegex *regexp.Regexp

var invalidationPrefixes = []string{
	"<h4>",
	"<th class=",
	"</table>",
	"<table",
	"<tr class=\"thead\">",
}

func init() {
	tableColRegex, err := regexp.Compile("<td class=\"c([0-9a-z]| )+\">")
	if err != nil {
		panic(err)
	}
	tdRegex = tableColRegex

	unitColRegex, err := regexp.Compile("(<td class=\"c([0-9a-z]| )+\">(dB|[a-zA-Z][bB]it/s)</td>| ms)")
	if err != nil {
		panic(err)
	}
	unitRegex = unitColRegex
}

func ParseConnectionInformation(html string) (*ConnectionInformation, error) {
	if len(html) == 0 {
		return nil, ParsingError
	}

	cleanLines := cleanupAndSplit(html)

	connectionMatrix := toTableMatrix(cleanLines)

	if len(connectionMatrix) < 16 {
		return nil, ParsingError
	}

	conInfo := ConnectionInformation{
		Upstream:   ConnectionTransmissionDirection{},
		Downstream: ConnectionTransmissionDirection{},
	}

	for index, row := range connectionMatrix {
		switch index {

		case 0: // DSLAM Max Datarate
			up, down, err := getUpDownValuesInt(row)
			if err != nil {
				return nil, err
			}

			conInfo.Downstream.MaximumDataRate = int(down)
			conInfo.Upstream.MaximumDataRate = int(up)

			break

		case 1: //DSLAM Min Datarate
			up, down, err := getUpDownValuesInt(row)
			if err != nil {
				return nil, err
			}

			conInfo.Downstream.MinimumDataRate = down
			conInfo.Upstream.MinimumDataRate = up

			break
		case 2: // Line Capacity
			up, down, err := getUpDownValuesInt(row)
			if err != nil {
				return nil, err
			}

			conInfo.Downstream.Capacity = down
			conInfo.Upstream.Capacity = up

			break
		case 3: // Current Data Rate
			up, down, err := getUpDownValuesInt(row)
			if err != nil {
				return nil, err
			}

			conInfo.Downstream.CurrentDataRate = down
			conInfo.Upstream.CurrentDataRate = up

			break
		case 4: // "Nahtlose Ratenadaption"
			up, down, err := getUpDownValuesBool(row)
			if err != nil {
				return nil, err
			}

			conInfo.Downstream.SeamlessRateAdjustment = down
			conInfo.Upstream.SeamlessRateAdjustment = up

			break
		case 5: // Latency
			up, down, err := getUpDownValuesInt(row)
			if err != nil {
				return nil, err
			}

			conInfo.Downstream.Latency = down
			conInfo.Upstream.Latency = up

			break
		case 6: // INP Value
			up, down, err := getUpDownValuesFloat(row)
			if err != nil {
				return nil, err
			}

			conInfo.Downstream.INPValue = down
			conInfo.Upstream.INPValue = up

			break
		case 7: // G.INP Value
			up, down, err := getUpDownValuesBool(row)
			if err != nil {
				return nil, err
			}

			conInfo.Downstream.GINP = down
			conInfo.Upstream.GINP = up

			break
		case 8: // SNR Value
			up, down, err := getUpDownValuesFloat(row)
			if err != nil {
				return nil, err
			}

			conInfo.Downstream.SNMargin = down
			conInfo.Upstream.SNMargin = up

			break
		case 9: // Bitswap Value
			up, down, err := getUpDownValuesBool(row)
			if err != nil {
				return nil, err
			}

			conInfo.Downstream.Bitswap = down
			conInfo.Upstream.Bitswap = up

			break
		case 10: // Line Attenuation Value
			up, down, err := getUpDownValuesFloat(row)
			if err != nil {
				return nil, err
			}

			conInfo.Downstream.LineAttenuation = down
			conInfo.Upstream.LineAttenuation = up

			break
		case 11: // Profile
			if len(row) != 2 {
				return nil, ParsingError
			}
			conInfo.Profile = row[1]

			break
		case 12: // G.Vector mode
			if len(row) != 3 {
				return nil, ParsingError
			}

			conInfo.Downstream.VectorMode = row[1]
			conInfo.Upstream.VectorMode = row[2]

			break
		case 13: // Carrier
			if len(row) != 3 {
				return nil, ParsingError
			}

			conInfo.Downstream.Carrier = row[1]
			conInfo.Upstream.Carrier = row[2]

			break
		case 14: // Downstream Errors
			errs, err := parseErrorFromTableLine(row)
			if err != nil {
				return nil, err
			}

			conInfo.Downstream.Errors = *errs
			break
		case 15: // Upstream Errors
			errs, err := parseErrorFromTableLine(row)
			if err != nil {
				return nil, err
			}

			conInfo.Upstream.Errors = *errs
			break
		}
	}

	return &conInfo, nil
}

func parseErrorFromTableLine(row []string) (*Errors, error) {
	if len(row) != 5 {
		return nil, ParsingError
	}

	es, err := strconv.ParseFloat(row[1], 64)
	if err != nil {
		return nil, ParsingError
	}

	ses, err := strconv.ParseFloat(row[2], 64)
	if err != nil {
		return nil, ParsingError
	}

	perMin, err := strconv.ParseFloat(row[3], 64)
	if err != nil {
		return nil, ParsingError
	}

	last15minES, err := strconv.ParseFloat(row[4], 64)
	if err != nil {
		return nil, ParsingError
	}

	return &Errors{
		SecondsWithErrors:     es,
		SecondsWithManyErrors: ses,
		ErrorsPerMinute:       perMin,
		ErrorsLast15Min:       last15minES,
	}, nil
}

func getUpDownValuesInt(row []string) (int, int, error) {
	if len(row) != 3 {
		return -1, -1, ParsingError
	}

	down, err := strconv.ParseInt(row[1], 10, 64)
	if row[1] == "fast" {
		down = 0
		err = nil
	}
	if err != nil {
		return -1, -1, ParsingError
	}

	up, err := strconv.ParseInt(row[2], 10, 64)
	if row[2] == "fast" {
		up = 0
		err = nil
	}
	if err != nil {
		return -1, -1, ParsingError
	}

	return int(up), int(down), nil
}

func getUpDownValuesFloat(row []string) (float64, float64, error) {
	if len(row) != 3 {
		return -1, -1, ParsingError
	}

	down, err := strconv.ParseFloat(row[1], 64)
	if err != nil {
		return -1, -1, ParsingError
	}

	up, err := strconv.ParseFloat(row[2], 64)
	if err != nil {
		return -1, -1, ParsingError
	}

	return up, down, nil
}

func getUpDownValuesBool(row []string) (bool, bool, error) {
	if len(row) != 3 {
		return false, false, ParsingError
	}

	trueregex, _ := regexp.Compile("(an|wahr|on|active|true)")

	down := trueregex.Match([]byte(strings.ToLower(row[1])))
	up := trueregex.Match([]byte(strings.ToLower(row[2])))

	return up, down, nil
}

func toTableMatrix(cleanLines []string) [][]string {
	connectionMatrix := make([][]string, 0)
	for _, line := range cleanLines {
		splitLine := strings.Split(line, "</td>")
		lineparts := make([]string, 0)
		for _, elem := range splitLine {
			elem = string(tdRegex.ReplaceAll([]byte(elem), []byte("")))

			if len(elem) > 0 && elem != "&nbsp;" {
				lineparts = append(lineparts, elem)
			}
		}

		if len(lineparts) > 0 {
			connectionMatrix = append(connectionMatrix, lineparts)
		}
	}
	return connectionMatrix
}

func cleanupAndSplit(html string) []string {
	lines := strings.Split(strings.Replace(html, "\n", "", -1), "</tr>")
	cleanLines := make([]string, 0)
	for _, v := range lines {
		v = string(unitRegex.ReplaceAll([]byte(v), []byte("")))
		v = strings.Replace(v, "<tr>", "", -1)

		if len(v) > 3 && !hasInvalidPrefix(v) {
			cleanLines = append(cleanLines, v)
		}
	}
	return cleanLines
}

func hasInvalidPrefix(line string) bool {
	for _, v := range invalidationPrefixes {
		if strings.HasPrefix(line, v) {
			return true
		}
	}

	return false
}
