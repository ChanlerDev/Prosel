export interface TopPage {
  path: string;
  refType?: string;
  refId?: string;
  views: number;
}

export interface TopReferer {
  referer: string;
  views: number;
}

export interface DeviceStat {
  deviceType: string;
  views: number;
}

export interface DailyView {
  date: string;
  views: number;
}

export interface AnalyticsOverview {
  todayViews: number;
  weekViews: number;
  monthViews: number;
  topPages: TopPage[];
  topReferers: TopReferer[];
  devices: DeviceStat[];
}

export interface PageViewRequest {
  path: string;
  refType?: 'post' | 'note' | 'page';
  refId?: string;
  referer?: string;
}
