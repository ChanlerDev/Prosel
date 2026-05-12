'use client';

import { NoteEditor } from '@/components/features/note/note-editor';
import { NoteStatusBadge } from '@/components/features/note/note-status-badge';
import { ApiErrorState, LoadingState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { useAdminNote, usePinNote, useUpdateNote } from '@/lib/notes/hooks';

export function AdminNoteEdit({ id }: { id: string }) {
  const note = useAdminNote(id);
  const update = useUpdateNote(id);
  const pin = usePinNote();

  if (note.isLoading) return <LoadingState />;
  if (note.isError || !note.data) return <ApiErrorState message="Unable to load note." />;

  return (
    <div className="grid gap-5">
      <div className="flex items-center gap-3">
        <NoteStatusBadge status={note.data.status} />
        <Button className="px-3 py-1 text-xs" disabled={pin.isPending} onClick={() => pin.mutate({ id: note.data.id, pinned: !note.data.pinnedAt })} type="button">
          {note.data.pinnedAt ? 'Unpin' : 'Pin'}
        </Button>
      </div>
      <NoteEditor error={update.isError ? update.error.message : undefined} isPending={update.isPending} note={note.data} onSubmit={(values) => update.mutate(values)} refId={note.data.id} />
    </div>
  );
}
