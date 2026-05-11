import { AdminPageEdit } from '@/components/features/page/admin-page-edit';
import { AdminShell } from '@/components/layout/admin-shell';

export default async function AdminPageEditPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">Edit page</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Update page content, template, and SEO metadata.</p>
      </div>
      <AdminPageEdit id={id} />
    </AdminShell>
  );
}
