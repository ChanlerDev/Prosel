'use client';

import { useQuery } from '@tanstack/react-query';

import { dashboardApi } from '@/lib/api/dashboard';

export const dashboardKeys = {
  overview: ['admin', 'dashboard', 'overview'] as const,
  activities: ['admin', 'activity-logs'] as const,
};

export function useDashboardOverview() {
  return useQuery({ queryKey: dashboardKeys.overview, queryFn: dashboardApi.overview });
}

export function useActivityLogs() {
  return useQuery({ queryKey: dashboardKeys.activities, queryFn: () => dashboardApi.activities(50) });
}
