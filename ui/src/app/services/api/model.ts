export interface StatusResponse {
  state: string;
  uptime: number;
}

export interface StatResponse {
  spectrum_count: number;
  stats: RepoStats;
  latest: LatestSpectrum;
}

export interface RepoStats {
  first_spectrum: number;
  latest_spectrum: number;
  total_count: number;
}

export interface LatestSpectrum {
  timestamp: number;
  date: DateKey;
}

export interface DateKey {
  day: string;
  month: string;
  year: string;
}
