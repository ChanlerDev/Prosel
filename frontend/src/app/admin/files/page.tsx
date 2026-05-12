import { AdminFiles } from '@/components/features/file/admin-files';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AdminFilesPage() {
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">Media library</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Upload images and reuse them in posts, notes, and pages.</p>
      </div>
      <AdminFiles />
    </AdminShell>
  );
}
