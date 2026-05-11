'use client';

import { ApiErrorState, LoadingState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { useRebuildSearchIndex, useSearchIndexStatus } from '@/lib/search/hooks';

export function AdminSearch() {
  const status = useSearchIndexStatus();
  const rebuild = useRebuildSearchIndex();
  const data = rebuild.data ?? status.data;

  return (
    <div className="grid gap-5">
      {status.isLoading ? <LoadingState /> : null}
      {status.isError ? <ApiErrorState message={status.error.message} /> : null}
      {rebuild.isError ? <ApiErrorState message={rebuild.error.message} /> : null}
      <Card className="grid gap-4 md:grid-cols-[1fr_auto] md:items-center">
        <div>
          <h2 className="text-xl font-semibold">Search index</h2>
          <p className="mt-2 text-sm text-[var(--muted-foreground)]">Rebuild the PostgreSQL full-text index from published posts, notes, and pages.</p>
        </div>
        <Button disabled={rebuild.isPending} onClick={() => rebuild.mutate()} type="button">{rebuild.isPending ? 'Rebuilding...' : 'Rebuild index'}</Button>
      </Card>
      {data ? (
        <div className="grid gap-4 md:grid-cols-4">
          <StatusCard label="Total" value={data.total} />
          <StatusCard label="Posts" value={data.posts} />
          <StatusCard label="Notes" value={data.notes} />
          <StatusCard label="Pages" value={data.pages} />
        </div>
      ) : null}
      {data?.updatedAt ? <p className="text-sm text-[var(--muted-foreground)]">Last updated {new Date(data.updatedAt).toLocaleString()}</p> : null}
    </div>
  );
}

function StatusCard({ label, value }: { label: string; value: number }) {
  return (
    <Card>
      <p className="text-sm text-[var(--muted-foreground)]">{label}</p>
      <p className="mt-2 text-3xl font-semibold">{value}</p>
    </Card>
  );
}
