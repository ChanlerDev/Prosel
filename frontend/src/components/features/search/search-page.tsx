'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

import { ApiErrorState, EmptyState, LoadingState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { useSearch } from '@/lib/search/hooks';
import type { SearchRefType, SearchResult } from '@/types/search';

const searchTypes: Array<{ value: SearchRefType | ''; label: string }> = [
  { value: '', label: 'All' },
  { value: 'post', label: 'Posts' },
  { value: 'note', label: 'Notes' },
  { value: 'page', label: 'Pages' },
];

export function SearchPage({ initialQuery, initialType, initialPage }: { initialQuery: string; initialType: SearchRefType | ''; initialPage: number }) {
  const router = useRouter();
  const [query, setQuery] = useState(initialQuery);
  const params = { q: initialQuery, type: initialType, page: initialPage, perPage: 10 };
  const search = useSearch(params);

  function submit(nextType = initialType, nextPage = 1) {
    const values = new URLSearchParams();
    if (query.trim()) values.set('q', query.trim());
    if (nextType) values.set('type', nextType);
    if (nextPage > 1) values.set('page', String(nextPage));
    router.push(`/search${values.toString() ? `?${values.toString()}` : ''}`);
  }

  return (
    <div className="grid gap-6">
      <form className="flex flex-col gap-3 md:flex-row" onSubmit={(event) => { event.preventDefault(); submit(); }}>
        <Input className="text-base" onChange={(event) => setQuery(event.target.value)} placeholder="Search posts, notes, and pages" value={query} />
        <Button type="submit">Search</Button>
      </form>
      <SearchTypeTabs current={initialType} onChange={(type) => submit(type)} />
      {!initialQuery ? <EmptyState title="Search the site" description="Enter a keyword to find published posts, notes, and pages." /> : null}
      {search.isLoading ? <LoadingState /> : null}
      {search.isError ? <ApiErrorState message={search.error.message} /> : null}
      {search.data && search.data.results.length === 0 ? <EmptyState title="No results" description="Try a different keyword or remove the type filter." /> : null}
      {search.data?.results.length ? <SearchResultList results={search.data.results} /> : null}
      {search.data && search.data.meta.totalPages > 1 ? (
        <div className="flex items-center justify-between text-sm text-[var(--muted-foreground)]">
          <span>Page {search.data.meta.page} of {search.data.meta.totalPages}</span>
          <div className="flex gap-2">
            <Button disabled={search.data.meta.page <= 1} onClick={() => submit(initialType, search.data.meta.page - 1)} type="button">Previous</Button>
            <Button disabled={search.data.meta.page >= search.data.meta.totalPages} onClick={() => submit(initialType, search.data.meta.page + 1)} type="button">Next</Button>
          </div>
        </div>
      ) : null}
    </div>
  );
}

function SearchTypeTabs({ current, onChange }: { current: SearchRefType | ''; onChange: (type: SearchRefType | '') => void }) {
  return (
    <div className="flex flex-wrap gap-2">
      {searchTypes.map((type) => (
        <button className={`rounded-full border px-4 py-2 text-sm ${current === type.value ? 'border-[var(--primary)] bg-[var(--primary)] text-[var(--primary-foreground)]' : 'border-[var(--border)] text-[var(--muted-foreground)]'}`} key={type.value || 'all'} onClick={() => onChange(type.value)} type="button">
          {type.label}
        </button>
      ))}
    </div>
  );
}

function SearchResultList({ results }: { results: SearchResult[] }) {
  return <div className="grid gap-4">{results.map((result) => <SearchResultItem key={`${result.refType}-${result.refId}`} result={result} />)}</div>;
}

function SearchResultItem({ result }: { result: SearchResult }) {
  return (
    <Card>
      <div className="mb-2 text-xs font-medium uppercase tracking-wide text-[var(--muted-foreground)]">{result.refType}</div>
      <Link className="text-xl font-semibold hover:text-[var(--primary)]" href={resultHref(result)}>{result.title}</Link>
      {result.excerpt ? <p className="mt-2 text-sm leading-6 text-[var(--muted-foreground)]">{result.excerpt}</p> : null}
    </Card>
  );
}

function resultHref(result: SearchResult) {
  if (result.refType === 'post') return `/posts/${result.slug}`;
  if (result.refType === 'note') return `/notes/${result.slug}`;
  return `/pages/${result.slug}`;
}
