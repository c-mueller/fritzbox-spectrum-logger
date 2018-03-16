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

import {Component, OnInit} from '@angular/core';
import {ApiService} from '../../services/api/api.service';
import {DateKey, SpectraKeyList} from '../../services/api/model';
import {sprintf} from 'sprintf-js';
import {Router} from '@angular/router';

@Component({
  selector: 'app-spectra',
  templateUrl: './spectra.component.html',
  styleUrls: ['./spectra.component.css']
})
export class SpectraComponent implements OnInit {

  public keys: SpectraKeyList = {keys: [], timestamp: 0};
  public selectedKey: DateKey = {day: '0', year: '0', month: '0'};
  public loading = false;
  public timestampList: number[][] = [];
  public expanded: boolean[] = [];

  constructor(private api: ApiService,
              private router: Router) {
  }

  selectDateKey(e: DateKey) {
    this.collapseAll();
    this.loading = true;
    this.selectedKey = e;
    this.api.getTimestampsForKey(e).subscribe(data => {
      this.clearTimestampList();
      for (const date of data.timestamps) {
        const d = new Date(date * 1000);
        this.timestampList[d.getHours()].push(date);
      }
      this.loading = false;
    });
  }

  clearTimestampList() {
    for (let i = 0; i < 24; i++) {
      this.timestampList[i] = [];
    }
  }

  collapseAll() {
    for (let i = 0; i < 24; i++) {
      this.expanded[i] = false;
    }
  }

  expandBlock(i: number) {
    this.collapseAll();
    this.expanded[i] = true;
  }

  fetchKeys() {
    this.api.getSpectrumKeys().subscribe(e => {
      this.keys = e;
    });
  }

  ngOnInit() {
    this.fetchKeys();
    this.collapseAll();
  }

  onDropdownClick(toggleValue: boolean, idx: number) {
    this.collapseAll();
    if (!toggleValue) {
      this.expandBlock(idx);
    }
  }

  showSpectrum(timestamp: number) {
    this.router.navigate(['/spectrum', timestamp]);
  }

  getHourQuarterForTimestamp(timestamp: number): number {
    const d = new Date(timestamp * 1000);
    return Math.floor(d.getMinutes() / (60 / 4));
  }

}
