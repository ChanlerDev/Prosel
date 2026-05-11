import { AdminFriends } from '@/components/features/page/admin-friends';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AdminFriendsPage() {
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">Friends</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Manage active, pending, and hidden friend links.</p>
      </div>
      <AdminFriends />
    </AdminShell>
  );
}
