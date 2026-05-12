import { request } from '@/lib/api/client';
import { readAuthState } from '@/lib/auth/store';
import type { AnalyticsOverview, DailyView, PageViewRequest, TopPage } from '@/types/analytics';

function accessToken() {
  const auth = readAuthState();
  if (!auth.accessToken) throw new Error('Authentication required');
  return auth.accessToken;
}

function queryString(params: object) {
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== '') search.set(key, String(value));
  }
  const value = search.toString();
  return value ? `?${value}` : '';
}

export const analyticsApi = {
  recordPageView: (body: PageViewRequest) => request<{ ok: boolean }>('/analytics/page-view', { method: 'POST', body: JSON.stringify(body), cache: 'no-store' }),
  admin: {
    overview: (range = '30d') => request<AnalyticsOverview>(`/admin/analytics/overview${queryString({ range })}`, { cache: 'no-store', accessToken: accessToken() }),
    topPages: (range = '30d', limit = 10) => request<TopPage[]>(`/admin/analytics/top-pages${queryString({ range, limit })}`, { cache: 'no-store', accessToken: accessToken() }),
    daily: (days = 30) => request<DailyView[]>(`/admin/analytics/daily${queryString({ days })}`, { cache: 'no-store', accessToken: accessToken() }),
  },
};
