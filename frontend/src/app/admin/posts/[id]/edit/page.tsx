import { AdminPostEdit } from '@/components/features/post/admin-post-edit';
import { AdminShell } from '@/components/layout/admin-shell';

export default async function AdminPostEditPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">Edit post</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Update Markdown, SEO metadata, and publishing state.</p>
      </div>
      <AdminPostEdit id={id} />
    </AdminShell>
  );
}
