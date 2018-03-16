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

export interface TimestampList {
  timestamps: number[];
  timestamp: number;
  requested_day: DateKey;
}

export interface NeighboursResponse {
  previous_timestamp: number;
  next_timestamp: number;
  request_timestamp: number;
}

export interface StatusResponse {
  state: string;
  uptime: number;
}

export interface InfoResponse {
  state: string;
  message: string;
}

export interface StatResponse {
  spectrum_count: number;
  stats: RepoStats;
  latest: LatestSpectrum;
}

export interface RepoStats {
  first_spectrum: number;
  latest_spectrum: number;
  total_count: number;
}

export interface LatestSpectrum {
  timestamp: number;
  date: DateKey;
}

export interface SpectraKeyList {
  timestamp: number;
  keys: DateKey[];
}

export interface DateKey {
  day: string;
  month: string;
  year: string;
}
