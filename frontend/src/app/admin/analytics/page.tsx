import { AdminAnalytics } from '@/components/features/analytics/admin-analytics';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AdminAnalyticsPage() {
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">Analytics</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Review traffic trends, top pages, sources, and device breakdown.</p>
      </div>
      <AdminAnalytics />
    </AdminShell>
  );
}
