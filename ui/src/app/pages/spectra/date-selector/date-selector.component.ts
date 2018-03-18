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
import {DateKey} from "../../../services/api/model";
import {sprintf} from "sprintf-js";

@Component({
  selector: 'app-date-selector',
  templateUrl: './date-selector.component.html',
  styleUrls: ['./date-selector.component.css']
})
export class DateSelectorComponent implements OnInit {

  @Input('dates') public validDates: DateKey[];
  @Output('selectedDate') public selectedDateEmitter = new EventEmitter<DateKey>();

  public selectedDateKey: DateKey = {month: "0", year: "0", day: "0"};

  constructor() {
  }

  ngOnInit() {
  }

  selectDateKey(e: DateKey) {
    this.selectedDateKey = e;
    this.selectedDateEmitter.emit(e);
  }

  getDateString(key: DateKey): string {
    const month = parseInt(key.month, 10);
    const year = parseInt(key.year, 10);
    const day = parseInt(key.day, 10);
    return sprintf('%02d.%02d.%04d', day, month, year);
  }

  keysEqual(a: DateKey, b: DateKey): boolean {
    return a != null && b != null && a.year === b.year && a.month === b.month && a.day === b.day;
  }

}
