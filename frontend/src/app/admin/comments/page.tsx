import { AdminComments } from '@/components/features/comment/admin-comments';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AdminCommentsPage() {
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">Comments</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Review visitor comments, moderate spam, and reply as the author.</p>
      </div>
      <AdminComments />
    </AdminShell>
  );
}
