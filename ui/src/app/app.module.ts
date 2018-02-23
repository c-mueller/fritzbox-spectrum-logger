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
  ],
  imports: [
    BrowserModule,
    RouterModule.forRoot(appRoutes),
    NgbModule.forRoot(),
    NgbAlertModule.forRoot()
  ],
  providers: [
    {
      provide: LocationStrategy,
      useClass: HashLocationStrategy
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
