import { AdminSubscribers } from '@/components/features/subscribe/admin-subscribers';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AdminSubscribersPage() {
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">Subscribers</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Review email subscribers and manually send new post notifications.</p>
      </div>
      <AdminSubscribers />
    </AdminShell>
  );
}
