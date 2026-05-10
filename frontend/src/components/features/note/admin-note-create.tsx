'use client';

import { NoteEditor } from '@/components/features/note/note-editor';
import { useCreateNote } from '@/lib/notes/hooks';

export function AdminNoteCreate() {
  const create = useCreateNote();
  return <NoteEditor error={create.isError ? create.error.message : undefined} isPending={create.isPending} onSubmit={(values) => create.mutate(values)} />;
}
