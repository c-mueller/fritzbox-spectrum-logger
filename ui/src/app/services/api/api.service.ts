import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {environment} from '../../../environments/environment';
import {StatResponse, StatusResponse} from './model';

@Injectable()
export class ApiService {

  private endpoint = environment.host;

  constructor(private http: HttpClient) {
  }

  getStatus() {
    return this.http.get<StatusResponse>(this.endpoint + '/api/status');
  }

  getStats() {
    return this.http.get<StatResponse>(this.endpoint + '/api/stats');
  }

}
