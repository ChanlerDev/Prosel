import { AdminNotesList } from '@/components/features/note/admin-notes-list';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AdminNotesPage() {
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">Notes</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Manage short-form updates, moods, weather, locations, and pinned notes.</p>
      </div>
      <AdminNotesList />
    </AdminShell>
  );
}
