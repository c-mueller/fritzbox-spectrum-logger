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

import {Component, EventEmitter, Input, OnChanges, OnInit, Output, SimpleChanges} from '@angular/core';
import {DateKey} from "../../../services/api/model";
import {NgbDateStruct} from "@ng-bootstrap/ng-bootstrap";

@Component({
  selector: 'app-date-selector',
  templateUrl: './date-selector.component.html',
  styleUrls: ['./date-selector.component.css']
})
export class DateSelectorComponent implements OnInit, OnChanges {
  @Input('dates') public validDates: DateKey[];
  @Output('selectedDate') public selectedDateEmitter = new EventEmitter<DateKey>();

  public selectedDateKey: DateKey = {month: "1", year: "2018", day: "1"};
  public minDate: NgbDateStruct = {month: 1, year: 2018, day: 1};
  public maxDate: NgbDateStruct = {month: 12, year: 2018, day: 31};

  public dateEnabledCallback: (value: NgbDateStruct, current: { year: number; month: number; }) => boolean = (value) => {
    let dk = this.toDateKey(value);

    for (let e of this.validDates) {
      if (this.keysEqual(e, dk)) {
        return false;
      }
    }

    return true;
  };

  onDateSelect(value: NgbDateStruct) {
    let dk = this.toDateKey(value);
    this.selectedDateEmitter.emit(dk);
  }

  ngOnInit() {
  }

  ngOnChanges(changes: SimpleChanges): void {
    this.minDate = this.toDateStruct(this.validDates[0]);
    this.maxDate = this.toDateStruct(this.validDates[this.validDates.length - 1])
  }

  toDateKey(value: NgbDateStruct): DateKey {
    return {year: "" + value.year, month: "" + value.month, day: "" + value.day};
  }

  toDateStruct(value: DateKey): NgbDateStruct {
    return {year: +value.year, month: +value.month, day: +value.day};
  }

  keysEqual(a: DateKey, b: DateKey): boolean {
    return a != null && b != null && a.year === b.year && a.month === b.month && a.day === b.day;
  }
}
