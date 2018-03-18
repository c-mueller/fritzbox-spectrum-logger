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

import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {sprintf} from "sprintf-js";

@Component({
  selector: 'app-spectra-hour-selector',
  templateUrl: './spectra-hour-selector.component.html',
  styleUrls: ['./spectra-hour-selector.component.css']
})
export class SpectraHourSelectorComponent implements OnInit {

  @Input('collapsed') public collapsed: boolean;
  @Input('timeIndex') public timeIndex: number;
  @Input('timestamps') public timestamps: number[];

  @Output('collapse') public collapse = new EventEmitter<boolean>();
  @Output('spectrumSelect') public spectrumSelect = new EventEmitter<number>();

  constructor() {
  }

  ngOnInit() {
    console.log(this.timeIndex)
  }

  onSpectrumSelect(timestamp: number) {
    this.spectrumSelect.emit(timestamp);
  }

  onToggleCollapse() {
    this.collapse.emit(!this.collapsed);
  }


  filterTimestampsByQuarter(quarter: number, timestamps: number[]): number[] {
    let data: number[] = [];
    for (let v of timestamps) {
      if (this.getHourQuarterForTimestamp(v) === quarter) {
        data.push(v);
      }
    }
    return data;
  }

  getHourQuarterForTimestamp(timestamp: number): number {
    const d = new Date(timestamp * 1000);
    return Math.floor(d.getMinutes() / (60 / 4));
  }

  formatHourIndexTime(i: number): string {
    return sprintf('%02d:00', i);
  }
}
