'use client';

import { useState } from 'react';

import { ApiErrorState, EmptyState, LoadingState } from '@/components/features/system/states';
import { Card } from '@/components/ui/card';
import { useAnalyticsOverview, useDailyViews } from '@/lib/analytics/hooks';

export function AdminAnalytics() {
  const [range, setRange] = useState('30d');
  const days = range === '7d' ? 7 : range === '90d' ? 90 : 30;
  const overview = useAnalyticsOverview(range);
  const daily = useDailyViews(days);

  if (overview.isLoading || daily.isLoading) return <LoadingState />;
  if (overview.isError || !overview.data) return <ApiErrorState message={overview.error?.message ?? 'Unable to load analytics.'} />;
  if (daily.isError) return <ApiErrorState message={daily.error.message} />;

  const maxDaily = Math.max(...(daily.data ?? []).map((item) => item.views), 1);

  return (
    <div className="grid gap-6">
      <Card className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h2 className="text-lg font-semibold">Traffic overview</h2>
          <p className="text-sm text-[var(--muted-foreground)]">Anonymous page views only. IP addresses are hashed before storage.</p>
        </div>
        <select className="rounded-xl border border-[var(--border)] bg-transparent px-3 py-2 text-sm" onChange={(event) => setRange(event.target.value)} value={range}>
          <option value="7d">Last 7 days</option>
          <option value="30d">Last 30 days</option>
          <option value="90d">Last 90 days</option>
        </select>
      </Card>

      <div className="grid gap-4 md:grid-cols-3">
        <MetricCard label="Today" value={overview.data.todayViews} />
        <MetricCard label="Last 7 days" value={overview.data.weekViews} />
        <MetricCard label="Last 30 days" value={overview.data.monthViews} />
      </div>

      <Card>
        <h2 className="mb-4 text-lg font-semibold">Daily views</h2>
        {daily.data && daily.data.length > 0 ? (
          <div className="flex h-48 items-end gap-2">
            {daily.data.map((item) => (
              <div className="flex flex-1 flex-col items-center gap-2" key={item.date} title={`${item.date}: ${item.views}`}>
                <div className="w-full rounded-t-xl bg-[var(--foreground)]/80" style={{ height: `${Math.max((item.views / maxDaily) * 100, 3)}%` }} />
                <span className="text-[10px] text-[var(--muted-foreground)]">{item.date.slice(5)}</span>
              </div>
            ))}
          </div>
        ) : <EmptyState title="No trend data" description="Page views will appear after visitors browse the site." />}
      </Card>

      <div className="grid gap-5 xl:grid-cols-3">
        <Card>
          <h2 className="mb-4 text-lg font-semibold">Top pages</h2>
          <StatList items={overview.data.topPages.map((page) => ({ label: page.path, value: page.views }))} empty="No page views yet." />
        </Card>
        <Card>
          <h2 className="mb-4 text-lg font-semibold">Sources</h2>
          <StatList items={overview.data.topReferers.map((source) => ({ label: source.referer, value: source.views }))} empty="No external sources yet." />
        </Card>
        <Card>
          <h2 className="mb-4 text-lg font-semibold">Devices</h2>
          <StatList items={overview.data.devices.map((device) => ({ label: device.deviceType, value: device.views }))} empty="No device data yet." />
        </Card>
      </div>
    </div>
  );
}

function MetricCard({ label, value }: { label: string; value: number }) {
  return (
    <Card>
      <p className="text-sm text-[var(--muted-foreground)]">{label}</p>
      <p className="mt-3 text-4xl font-semibold">{value}</p>
    </Card>
  );
}

function StatList({ items, empty }: { items: { label: string; value: number }[]; empty: string }) {
  if (items.length === 0) return <p className="text-sm text-[var(--muted-foreground)]">{empty}</p>;
  return (
    <div className="grid gap-3">
      {items.map((item) => (
        <div className="flex items-center justify-between gap-4 text-sm" key={item.label}>
          <span className="truncate text-[var(--muted-foreground)]">{item.label}</span>
          <span className="font-semibold">{item.value}</span>
        </div>
      ))}
    </div>
  );
}
