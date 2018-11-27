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
	"reflect"
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

	conInfo := ConnectionInformation{
		Upstream:   ConnectionTransmissionDirection{},
		Downstream: ConnectionTransmissionDirection{},
		LineLength: -1,
	}

	for index, row := range connectionMatrix {

		key := row[0]

		t := reflect.TypeOf(conInfo.Downstream)

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)

			value, present := f.Tag.Lookup("fbname")
			if present && value == key {
				downVal := reflect.ValueOf(&conInfo.Downstream).Elem().FieldByName(f.Name)
				upVal := reflect.ValueOf(&conInfo.Upstream).Elem().FieldByName(f.Name)

				downStr := row[1]
				upStr := row[2]

				switch f.Type.Kind() {
				case reflect.Int:
					downVal.SetInt(int64(parseInteger(downStr)))
					upVal.SetInt(int64(parseInteger(upStr)))
				case reflect.String:
					downVal.SetString(downStr)
					upVal.SetString(upStr)
				case reflect.Bool:
					downVal.SetBool(parseBoolean(downStr))
					upVal.SetBool(parseBoolean(upStr))
				case reflect.Float64:
					downVal.SetFloat(parseFloat(downStr))
					upVal.SetFloat(parseFloat(upStr))
				}
			}

		}

		if key == "Profil" {
			conInfo.Profile = row[1]
		} else if key == "ungefähre Leitungslänge" {
			conInfo.LineLength = parseInteger(row[2])
		}

		switch index {
		case len(connectionMatrix) - 2: // Downstream Errors
			errs, err := parseErrorFromTableLine(row)
			if err != nil {
				return nil, err
			}

			conInfo.Downstream.Errors = *errs
			break
		case len(connectionMatrix) - 2: // Upstream Errors
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

func parseFloat(val string) float64 {
	v, err := strconv.ParseFloat(val, 64)

	if err != nil {
		return -1
	}
	return v
}

func parseBoolean(val string) bool {
	trueregex, _ := regexp.Compile("(an|wahr|on|active|true)")

	return trueregex.Match([]byte(val))
}

func parseInteger(val string) int {
	if val == "fast" {
		return 0
	} else if val == "-" {
		return 0
	}

	value, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return -1
	}

	return int(value)
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
