import { AdminActivity } from '@/components/features/dashboard/admin-activity';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AdminActivityPage() {
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">Activity</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Recent backend activity records.</p>
      </div>
      <AdminActivity />
    </AdminShell>
  );
}
