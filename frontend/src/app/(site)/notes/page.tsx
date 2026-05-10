import { NoteTimeline } from '@/components/features/note/note-timeline';
import { EmptyState } from '@/components/features/system/states';
import { SiteContainer } from '@/components/layout/site-container';
import { SiteFooter } from '@/components/layout/site-footer';
import { SiteHeader } from '@/components/layout/site-header';
import { notesApi } from '@/lib/api/notes';

export const revalidate = 60;

export default async function NotesPage({ searchParams }: { searchParams: Promise<{ page?: string }> }) {
  const { page } = await searchParams;
  const currentPage = Number(page ?? '1') || 1;
  const result = await notesApi.list({ page: currentPage, perPage: 20 }).catch(() => ({ notes: [], meta: { page: 1, perPage: 20, total: 0, totalPages: 0 } }));

  return (
    <>
      <SiteHeader />
      <main className="py-12">
        <SiteContainer>
          <div className="mx-auto mb-10 max-w-3xl">
            <h1 className="text-5xl font-semibold tracking-tight">Notes</h1>
            <p className="mt-4 text-lg text-[var(--muted-foreground)]">Short updates, field notes, and small moments from the blog.</p>
          </div>
          <div className="mx-auto max-w-3xl">
            {result.notes.length > 0 ? <NoteTimeline notes={result.notes} /> : <EmptyState title="No notes yet" description="Short updates will appear here once published." />}
          </div>
        </SiteContainer>
      </main>
      <SiteFooter />
    </>
  );
}
