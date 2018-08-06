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
import {FSLConfiguration} from "../../services/api/model";
import {ApiService} from "../../services/api/api.service";

@Component({
  selector: 'app-config',
  templateUrl: './config.component.html',
  styleUrls: ['./config.component.css']
})
export class ConfigComponent implements OnInit {

  currentConfig: FSLConfiguration = null;
  loading: boolean = true;
  error: boolean = false;

  constructor(private api: ApiService) {
  }

  ngOnInit() {
    this.api.getConfiguration().subscribe(e => {
      this.currentConfig = e;
      this.loading = false;
    }, err => {
      this.loading = false;
      this.error = true;
    })
  }

}
