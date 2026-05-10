import type { Metadata } from 'next';
import { notFound } from 'next/navigation';

import { CommentSection } from '@/components/features/comment/comment-section';
import { NoteContent } from '@/components/features/note/note-content';
import { Badge } from '@/components/ui/badge';
import { SiteContainer } from '@/components/layout/site-container';
import { SiteFooter } from '@/components/layout/site-footer';
import { SiteHeader } from '@/components/layout/site-header';
import { notesApi } from '@/lib/api/notes';

export const revalidate = 60;

export async function generateMetadata({ params }: { params: Promise<{ slug: string }> }): Promise<Metadata> {
  const { slug } = await params;
  const note = await notesApi.detail(slug).catch(() => null);
  if (!note) return {};
  return { title: note.title || `Note ${formatDate(note.publishedAt ?? note.createdAt)}`, description: note.contentText.slice(0, 160) };
}

export default async function NoteDetailPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;
  const note = await notesApi.detail(slug).catch(() => null);
  if (!note) notFound();

  return (
    <>
      <SiteHeader />
      <main className="py-12">
        <SiteContainer>
          <article className="mx-auto max-w-3xl">
            <div className="mb-8">
              <Badge>{note.viewCount} views</Badge>
              <h1 className="mt-5 text-5xl font-semibold tracking-tight">{note.title || 'Note'}</h1>
              <p className="mt-4 text-sm text-[var(--muted-foreground)]">Published {formatDate(note.publishedAt ?? note.createdAt)}{note.mood ? ` · ${note.mood}` : ''}{note.weather ? ` · ${note.weather}` : ''}{note.location ? ` · ${note.location}` : ''}</p>
            </div>
            <NoteContent note={note} />
          </article>
          <CommentSection refId={note.id} refType="note" />
        </SiteContainer>
      </main>
      <SiteFooter />
    </>
  );
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat('en', { dateStyle: 'medium' }).format(new Date(value));
}
