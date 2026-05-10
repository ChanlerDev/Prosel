import { AdminNoteEdit } from '@/components/features/note/admin-note-edit';
import { AdminShell } from '@/components/layout/admin-shell';

export default async function AdminNoteEditPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">Edit note</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Update timeline content and publishing state.</p>
      </div>
      <AdminNoteEdit id={id} />
    </AdminShell>
  );
}
