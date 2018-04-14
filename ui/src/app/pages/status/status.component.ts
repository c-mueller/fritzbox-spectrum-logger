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
import {StatResponse, StatusResponse} from '../../services/api/model';
import {sprintf} from 'sprintf-js';

@Component({
  selector: 'app-status',
  templateUrl: './status.component.html',
  styleUrls: ['./status.component.css']
})
export class StatusComponent implements OnInit {

  status: StatusResponse = {state: '', uptime: -1};
  stats: StatResponse = {latest: null, spectrum_count: -1, stats: null};

  numberWithDots = (x) => {
    return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".");
  };

  formatDuration = (x: number) => {
    const hours = Math.floor(x / 3600);
    const minutes = Math.floor((x - (hours * 3600)) / 60);
    const seconds = x - (hours * 3600) - (minutes * 60);

    return sprintf("%02d:%02d:%02d", hours, minutes, seconds);
  };

  constructor(private api: ApiService) {
  }

  toggleLogging() {
    console.log('Toggle Logging');
    if (this.status.state === 'LOGGING') {
      this.api.stopLogging().subscribe(e => {
        this.updateStatus();
      });
    } else {
      this.api.startLogging().subscribe(e => {
        this.updateStatus();
      });
    }
  }

  refresh() {
    this.updateStatus();
    this.updateStats();
  }

  ngOnInit() {
    this.refresh();
  }

  private updateStats() {
    this.api.getStats().subscribe(e => {
      this.stats = e;
    });
  }

  private updateStatus() {
    this.api.getStatus().subscribe(e => {
      this.status = e;
    });
  }

  toTime(val: number): string {
    return new Date(val * 1000).toLocaleString();
  }

}
