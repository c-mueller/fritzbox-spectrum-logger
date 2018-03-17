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

import {Component, Output, Input, OnInit, EventEmitter} from '@angular/core';
import {sprintf} from "sprintf-js";

@Component({
  selector: 'app-quarter-timestamp-selector',
  templateUrl: './quarter-timestamp-selector.component.html',
  styleUrls: ['./quarter-timestamp-selector.component.css']
})
export class QuarterTimestampSelectorComponent implements OnInit {

  @Input('quarterNumber') qIdx: number;
  @Input('hourNumber') hIdx: number;
  @Input('timestamps') timestamps: number;

  @Output('onSelect') selectEmitter = new EventEmitter<number>();

  constructor() {
  }

  onItemSelected(timestamp: number) {
    this.selectEmitter.emit(timestamp);
  }

  ngOnInit() {
  }


  formatSpectrumButtonTime(timestamp: number): string {
    const date = new Date(timestamp * 1000);
    return sprintf('%02d:%02d:%02d', date.getHours(), date.getMinutes(), date.getSeconds());
  }

  getTimeRangeStringForQuarter(quarter: number, hour: number) {
    let toHour = hour;
    let toMinute = (quarter + 1) * 15;
    if (toMinute >= 60) {
      toMinute = 0;
      toHour = ((toHour + 1) % 24);
    }
    return sprintf('%02d:%02d to %02d:%02d', hour, quarter * 15, toHour, toMinute);
  }

}
