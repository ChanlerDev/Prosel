import { AdminPageCreate } from '@/components/features/page/admin-page-create';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AdminPageCreatePage() {
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">New page</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Write a custom Markdown page and publish it to the public site.</p>
      </div>
      <AdminPageCreate />
    </AdminShell>
  );
}
