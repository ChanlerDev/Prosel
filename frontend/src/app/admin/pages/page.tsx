import { AdminPagesList } from '@/components/features/page/admin-pages-list';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AdminPagesPage() {
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">Pages</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Manage custom pages such as about, friends, and projects.</p>
      </div>
      <AdminPagesList />
    </AdminShell>
  );
}
