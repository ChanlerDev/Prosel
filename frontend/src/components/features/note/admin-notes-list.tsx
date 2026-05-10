'use client';

import Link from 'next/link';
import { useState } from 'react';

import { NoteStatusBadge } from '@/components/features/note/note-status-badge';
import { ApiErrorState, EmptyState, LoadingState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { useAdminNotes, useDeleteNote, usePinNote } from '@/lib/notes/hooks';
import type { NoteStatus } from '@/types/note';

export function AdminNotesList() {
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState<NoteStatus | ''>('');
  const notes = useAdminNotes({ search, status });
  const pin = usePinNote();
  const remove = useDeleteNote();

  return (
    <div className="grid gap-5">
      <Card className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div className="flex gap-3">
          <Input onChange={(event) => setSearch(event.target.value)} placeholder="Search notes" value={search} />
          <select className="rounded-xl border border-[var(--border)] bg-transparent px-3 py-2 text-sm" onChange={(event) => setStatus(event.target.value as NoteStatus | '')} value={status}>
            <option value="">All</option>
            <option value="published">Published</option>
            <option value="draft">Draft</option>
            <option value="private">Private</option>
            <option value="archived">Archived</option>
          </select>
        </div>
        <Link className="rounded-full bg-[var(--primary)] px-4 py-2 text-sm font-medium text-[var(--primary-foreground)]" href="/admin/notes/new">
          New note
        </Link>
      </Card>
      {notes.isLoading ? <LoadingState /> : null}
      {notes.isError ? <ApiErrorState message={notes.error.message} /> : null}
      {notes.data && notes.data.notes.length === 0 ? <EmptyState title="No notes" description="Create short updates for your public timeline." /> : null}
      {notes.data?.notes.map((note) => (
        <Card className="grid gap-3 md:grid-cols-[1fr_auto] md:items-center" key={note.id}>
          <div>
            <div className="flex items-center gap-2">
              <NoteStatusBadge status={note.status} />
              {note.pinnedAt ? <span className="text-xs text-[var(--muted-foreground)]">Pinned</span> : null}
            </div>
            <Link className="mt-2 block text-xl font-semibold hover:text-[var(--primary)]" href={`/admin/notes/${note.id}/edit`}>
              {note.title || note.contentText.slice(0, 80) || 'Untitled note'}
            </Link>
            <p className="mt-1 text-sm text-[var(--muted-foreground)]">/{note.slug} · {note.viewCount} views · {formatDate(note.updatedAt)}</p>
          </div>
          <div className="flex flex-wrap gap-2">
            <Button className="px-3 py-1 text-xs" disabled={pin.isPending} onClick={() => pin.mutate({ id: note.id, pinned: !note.pinnedAt })} type="button">
              {note.pinnedAt ? 'Unpin' : 'Pin'}
            </Button>
            <Link className="rounded-full border border-[var(--border)] px-3 py-1 text-xs" href={`/admin/notes/${note.id}/edit`}>Edit</Link>
            <Button className="bg-red-600 px-3 py-1 text-xs text-white" disabled={remove.isPending} onClick={() => remove.mutate(note.id)} type="button">Delete</Button>
          </div>
        </Card>
      ))}
    </div>
  );
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat('en', { dateStyle: 'medium' }).format(new Date(value));
}
