import { SearchPage } from '@/components/features/search/search-page';
import { SiteContainer } from '@/components/layout/site-container';
import { SiteFooter } from '@/components/layout/site-footer';
import { SiteHeader } from '@/components/layout/site-header';
import type { SearchRefType } from '@/types/search';

export const dynamic = 'force-dynamic';

export default async function PublicSearchPage({ searchParams }: { searchParams: Promise<{ q?: string; type?: string; page?: string }> }) {
  const params = await searchParams;
  const type = normalizeType(params.type);
  const page = Number(params.page ?? '1');

  return (
    <>
      <SiteHeader />
      <main className="py-12">
        <SiteContainer>
          <div className="mb-8">
            <h1 className="text-4xl font-semibold tracking-tight">Search</h1>
            <p className="mt-3 text-[var(--muted-foreground)]">Find published posts, notes, and pages.</p>
          </div>
          <SearchPage initialPage={Number.isFinite(page) && page > 0 ? page : 1} initialQuery={params.q ?? ''} initialType={type} />
        </SiteContainer>
      </main>
      <SiteFooter />
    </>
  );
}

function normalizeType(value?: string): SearchRefType | '' {
  return value === 'post' || value === 'note' || value === 'page' ? value : '';
}
