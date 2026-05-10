import { AdminDashboard } from '@/components/features/dashboard/admin-dashboard';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AdminHomePage() {
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">Dashboard</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Overview of your writing workflow, taxonomy, and recent activity.</p>
      </div>
      <AdminDashboard />
    </AdminShell>
  );
}
