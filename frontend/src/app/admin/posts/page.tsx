import { AdminPostsList } from '@/components/features/post/admin-posts-list';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AdminPostsPage() {
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">Posts</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Create drafts, publish articles, and manage blog content.</p>
      </div>
      <AdminPostsList />
    </AdminShell>
  );
}
