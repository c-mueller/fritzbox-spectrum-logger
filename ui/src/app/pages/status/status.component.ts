import {Component, OnInit} from '@angular/core';
import {ApiService} from '../../services/api/api.service';
import {StatResponse, StatusResponse} from '../../services/api/model';

@Component({
  selector: 'app-status',
  templateUrl: './status.component.html',
  styleUrls: ['./status.component.css']
})
export class StatusComponent implements OnInit {

  status: StatusResponse = {state: '', uptime: -1};
  stats: StatResponse = {latest: null, spectrum_count: -1, stats: null};

  constructor(private api: ApiService) {
  }

  ngOnInit() {
    this.api.getStatus().subscribe(e => {
      this.status = e;
    });
    this.api.getStats().subscribe(e => {
      this.stats = e;
    });
  }

  toTime(val: number): string {
    return new Date(val * 1000).toLocaleString();
  }

}
