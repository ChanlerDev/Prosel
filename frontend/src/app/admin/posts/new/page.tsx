import { AdminPostCreate } from '@/components/features/post/admin-post-create';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AdminPostCreatePage() {
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">New post</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Write a Markdown draft and publish it when ready.</p>
      </div>
      <AdminPostCreate />
    </AdminShell>
  );
}
