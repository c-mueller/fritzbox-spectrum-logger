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

import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {ApiService} from '../../services/api/api.service';
import {environment} from '../../../environments/environment';
import {sprintf} from 'sprintf-js';

@Component({
  selector: 'app-spectrum-view',
  templateUrl: './spectrum-view.component.html',
  styleUrls: ['./spectrum-view.component.css']
})
export class SpectrumViewComponent implements OnInit, OnDestroy {

  private subscription: any;
  public timestamp: number;
  public initialLoading = true;
  public loading = true;
  public data: string;

  constructor(private route: ActivatedRoute,
              private router: Router,
              private api: ApiService) {
  }

  ngOnInit() {
    this.subscription = this.route.params.subscribe(params => {
      const param = params['timestamp'];
      if (param == null || param === undefined) {
        this.router.navigate(['/404']);
      }
      const timestamp = parseInt(param, 10);
      if (timestamp === -1 || Number.isNaN(timestamp) || timestamp === 0) {
        this.router.navigate(['/404']);
      }
      this.timestamp = timestamp;
      this.requestImage();
    });
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  requestImage() {
    this.api.getSpectrumImage(this.timestamp).subscribe(image => {
        const reader = new FileReader();
        reader.addEventListener('load', () => {
          this.data = reader.result;
        }, false);

        if (image) {
          reader.readAsDataURL(image);
          this.initialLoading = false;
          this.loading = false;
        } else {
          this.router.navigate(['/404']);
        }
      },
      err => {
        this.router.navigate(['/404']);
      });
  }

  navigateTo(timestamp: number) {
    this.loading = true;
    this.router.navigate(['/spectrum', timestamp]);
  }

  getImageUrl(): string {
    return sprintf('%s/spectrum/%d/img', environment.host, this.timestamp);
  }

  formatTimestamp(timestamp: number): string {
    if (timestamp === -1) {
      return 'None';
    }
    const date = new Date(timestamp * 1000);
    return sprintf('%02d:%02d:%02d', date.getHours(), date.getMinutes(), date.getSeconds());
  }
}
