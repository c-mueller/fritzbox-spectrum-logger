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

import {Component, EventEmitter, Input, OnChanges, OnInit, Output} from '@angular/core';
import {ApiService} from "../../../services/api/api.service";
import {NeighboursResponse} from "../../../services/api/model";
import {sprintf} from "sprintf-js";

@Component({
  selector: 'app-spectrum-pagination',
  templateUrl: './spectrum-pagination.component.html',
  styleUrls: ['./spectrum-pagination.component.css']
})
export class SpectrumPaginationComponent implements OnInit, OnChanges {

  @Input('timestamp') public timestamp: number;
  @Output('onTimestampChange') public timestampEmitter = new EventEmitter<number>();

  public neighbours: NeighboursResponse = {previous_timestamp: -1, next_timestamp: -1, request_timestamp: -1};

  constructor(private api: ApiService) {
  }

  ngOnInit() {
    this.fetchNeighbours();
  }

  ngOnChanges() {
    this.fetchNeighbours();
  }

  emitTimestampChangeEvent(timestamp: number) {
    this.timestampEmitter.emit(timestamp);
  }

  fetchNeighbours() {
    this.api.getSpectrumNeighbours(this.timestamp).subscribe(e => {
      this.neighbours = e;
    })
  }

  formatTimestamp(timestamp: number): string {
    if (timestamp === -1) {
      return 'None';
    }
    const date = new Date(timestamp * 1000);
    return sprintf('%02d:%02d:%02d', date.getHours(), date.getMinutes(), date.getSeconds());
  }

}
