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

import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {environment} from '../../../environments/environment';
import {DateKey, InfoResponse, NeighboursResponse, SpectraKeyList, StatResponse, StatusResponse, TimestampList} from './model';

@Injectable()
export class ApiService {

  private endpoint = environment.host;

  constructor(private http: HttpClient) {
  }

  getSpectrumKeys() {
    return this.http.get<SpectraKeyList>(this.endpoint + '/spectra');
  }

  getTimestampsForKey(key: DateKey) {
    return this.http.get<TimestampList>(this.endpoint + '/spectra/' + key.year + '/' + key.month + '/' + key.day);
  }

  getSpectrumNeighbours(timestamp: number) {
    return this.http.get<NeighboursResponse>(this.endpoint + '/spectrum/' + timestamp + '/neighbours');
  }

  getSpectrumImage(timestamp: number) {
    return this.http.get(this.endpoint + '/spectrum/' + timestamp + '/img', {responseType: 'blob'});
  }

  startLogging() {
    return this.http.post<InfoResponse>(this.endpoint + '/control/start', null);
  }

  stopLogging() {
    return this.http.post<InfoResponse>(this.endpoint + '/control/stop', null);
  }

  getStatus() {
    return this.http.get<StatusResponse>(this.endpoint + '/status');
  }

  getStats() {
    return this.http.get<StatResponse>(this.endpoint + '/stats');
  }

}
