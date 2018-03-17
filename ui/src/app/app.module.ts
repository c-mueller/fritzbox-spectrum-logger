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

import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import {AppComponent} from './app.component';
import {NgbAlertModule, NgbModule} from '@ng-bootstrap/ng-bootstrap';
import {RouterModule, Routes} from '@angular/router';
import {NotFoundComponent} from './pages/not-found/not-found.component';
import {StatusComponent} from './pages/status/status.component';
import {HashLocationStrategy, LocationStrategy} from '@angular/common';
import {AboutComponent} from './pages/about/about.component';
import {ConfigComponent} from './pages/config/config.component';
import {SpectraComponent} from './pages/spectra/spectra.component';
import {HttpClientModule} from '@angular/common/http';
import {ApiService} from './services/api/api.service';
import {SpectrumViewComponent} from './pages/spectrum-view/spectrum-view.component';
import {DateSelectorComponent} from './pages/spectra/date-selector/date-selector.component';
import {SpectraHourSelectorComponent} from './pages/spectra/spectra-hour-selector/spectra-hour-selector.component';
import {QuarterTimestampSelectorComponent} from './pages/spectra/quarter-timestamp-selector/quarter-timestamp-selector.component';
import {SpectrumPaginationComponent} from './pages/spectrum-view/spectrum-pagination/spectrum-pagination.component';


const appRoutes: Routes = [
  {
    path: 'status',
    component: StatusComponent
  },
  {
    path: 'spectra',
    component: SpectraComponent
  },
  {
    path: 'spectrum/:timestamp',
    component: SpectrumViewComponent
  },
  {
    path: 'config',
    component: ConfigComponent
  },
  {
    path: 'about',
    component: AboutComponent
  },
  {
    path: '',
    redirectTo: '/status',
    pathMatch: 'full'
  },
  {
    path: '**',
    component: NotFoundComponent
  }
];

@NgModule({
  declarations: [
    AppComponent,
    AboutComponent,
    ConfigComponent,
    SpectraComponent,
    StatusComponent,
    NotFoundComponent,
    SpectrumViewComponent,
    DateSelectorComponent,
    SpectraHourSelectorComponent,
    QuarterTimestampSelectorComponent,
    SpectrumPaginationComponent,
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    BrowserModule,
    RouterModule.forRoot(appRoutes),
    NgbModule.forRoot(),
    NgbAlertModule.forRoot()
  ],
  providers: [
    {
      provide: LocationStrategy,
      useClass: HashLocationStrategy
    },
    ApiService,
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
