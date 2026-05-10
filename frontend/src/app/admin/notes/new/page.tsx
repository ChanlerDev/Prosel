import { AdminNoteCreate } from '@/components/features/note/admin-note-create';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AdminNoteCreatePage() {
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">New note</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Write a short Markdown update for the notes timeline.</p>
      </div>
      <AdminNoteCreate />
    </AdminShell>
  );
}
