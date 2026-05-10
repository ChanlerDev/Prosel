import Link from 'next/link';

import { Badge } from '@/components/ui/badge';
import { Card } from '@/components/ui/card';
import type { Note } from '@/types/note';

export function NoteTimeline({ notes }: { notes: Note[] }) {
  return (
    <div className="grid gap-5">
      {notes.map((note) => (
        <NoteCard key={note.id} note={note} />
      ))}
    </div>
  );
}

export function NoteCard({ note }: { note: Note }) {
  return (
    <Card className="relative overflow-hidden">
      {note.pinnedAt ? <PinnedNoteBadge /> : null}
      <div className="flex flex-wrap items-center gap-2 text-sm text-[var(--muted-foreground)]">
        <span>{formatDate(note.publishedAt ?? note.createdAt)}</span>
        {note.mood ? <span>· {note.mood}</span> : null}
        {note.weather ? <span>· {note.weather}</span> : null}
        {note.location ? <span>· {note.location}</span> : null}
      </div>
      {note.title ? (
        <Link className="mt-3 block text-2xl font-semibold tracking-tight hover:text-[var(--primary)]" href={`/notes/${note.slug}`}>
          {note.title}
        </Link>
      ) : null}
      <p className="mt-4 whitespace-pre-wrap text-[var(--foreground)] leading-8">{note.contentText}</p>
      <div className="mt-5 flex items-center gap-3 text-xs text-[var(--muted-foreground)]">
        <span>{note.viewCount} views</span>
        <span>{note.commentCount} comments</span>
        <Link className="font-medium text-[var(--primary)]" href={`/notes/${note.slug}`}>Open</Link>
      </div>
    </Card>
  );
}

function PinnedNoteBadge() {
  return <Badge className="absolute right-5 top-5">Pinned</Badge>;
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat('en', { dateStyle: 'medium' }).format(new Date(value));
}
