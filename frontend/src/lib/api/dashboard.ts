import { request } from '@/lib/api/client';
import { readAuthState } from '@/lib/auth/store';
import type { ActivityLog, DashboardOverview, DashboardStats } from '@/types/dashboard';

function accessToken() {
  const auth = readAuthState();
  if (!auth.accessToken) throw new Error('Authentication required');
  return auth.accessToken;
}

export const dashboardApi = {
  overview: () => request<DashboardOverview>('/admin/dashboard/overview', { cache: 'no-store', accessToken: accessToken() }),
  stats: () => request<DashboardStats>('/admin/dashboard/stats', { cache: 'no-store', accessToken: accessToken() }),
  activities: (limit = 20) => request<ActivityLog[]>(`/admin/activity-logs?limit=${limit}`, { cache: 'no-store', accessToken: accessToken() }),
};
