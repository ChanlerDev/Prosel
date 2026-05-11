import { AdminSearch } from '@/components/features/search/admin-search';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AdminSearchPage() {
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">Search</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Monitor and rebuild the site search index.</p>
      </div>
      <AdminSearch />
    </AdminShell>
  );
}
