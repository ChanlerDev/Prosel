'use client';

import { ActivityTimeline } from '@/components/features/dashboard/activity-timeline';
import { QuickActionGrid } from '@/components/features/dashboard/quick-action-grid';
import { RecentPostTable } from '@/components/features/dashboard/recent-post-table';
import { StatCard } from '@/components/features/dashboard/stat-card';
import { ApiErrorState, LoadingState } from '@/components/features/system/states';
import { useDashboardOverview } from '@/lib/dashboard/hooks';

export function AdminDashboard() {
  const overview = useDashboardOverview();

  if (overview.isLoading) return <LoadingState />;
  if (overview.isError || !overview.data) return <ApiErrorState message={overview.error?.message ?? 'Unable to load dashboard.'} />;

  const { stats, recentPosts, activities } = overview.data;
  return (
    <div className="grid gap-8">
      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <StatCard detail={`${stats.publishedPosts} published · ${stats.draftPosts} drafts`} title="Posts" value={stats.totalPosts} />
        <StatCard detail={`${stats.todayViews} today · across all pages`} title="Views" value={stats.totalViews} />
        <StatCard detail={`${stats.tags} tags · ${stats.topics} topics`} title="Taxonomy" value={stats.categories} />
        <StatCard detail="Comment module pending" title="Pending comments" value={stats.pendingComments} />
      </div>
      <section>
        <h2 className="mb-4 text-xl font-semibold">Quick actions</h2>
        <QuickActionGrid />
      </section>
      <div className="grid gap-5 xl:grid-cols-[1.4fr_1fr]">
        <RecentPostTable posts={recentPosts} />
        <ActivityTimeline activities={activities} />
      </div>
    </div>
  );
}
