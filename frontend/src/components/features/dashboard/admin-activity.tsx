'use client';

import { ActivityTimeline } from '@/components/features/dashboard/activity-timeline';
import { ApiErrorState, LoadingState } from '@/components/features/system/states';
import { useActivityLogs } from '@/lib/dashboard/hooks';

export function AdminActivity() {
  const activities = useActivityLogs();
  if (activities.isLoading) return <LoadingState />;
  if (activities.isError || !activities.data) return <ApiErrorState message={activities.error?.message ?? 'Unable to load activity.'} />;
  return <ActivityTimeline activities={activities.data} />;
}
