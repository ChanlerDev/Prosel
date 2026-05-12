'use client';

import { useQuery } from '@tanstack/react-query';

import { analyticsApi } from '@/lib/api/analytics';

export const analyticsKeys = {
  overview: (range: string) => ['admin', 'analytics', 'overview', range] as const,
  daily: (days: number) => ['admin', 'analytics', 'daily', days] as const,
};

export function useAnalyticsOverview(range: string) {
  return useQuery({ queryKey: analyticsKeys.overview(range), queryFn: () => analyticsApi.admin.overview(range) });
}

export function useDailyViews(days: number) {
  return useQuery({ queryKey: analyticsKeys.daily(days), queryFn: () => analyticsApi.admin.daily(days) });
}
