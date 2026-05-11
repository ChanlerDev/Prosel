'use client';

import Link from 'next/link';
import { useState } from 'react';

import { PageStatusBadge } from '@/components/features/page/page-status-badge';
import { ApiErrorState, EmptyState, LoadingState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { useAdminPages, useDeletePage } from '@/lib/pages/hooks';
import type { PageStatus } from '@/types/page';

export function AdminPagesList() {
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState<PageStatus | ''>('');
  const pages = useAdminPages({ search, status });
  const remove = useDeletePage();

  return (
    <div className="grid gap-5">
      <Card className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div className="flex gap-3">
          <Input onChange={(event) => setSearch(event.target.value)} placeholder="Search pages" value={search} />
          <select className="rounded-xl border border-[var(--border)] bg-transparent px-3 py-2 text-sm" onChange={(event) => setStatus(event.target.value as PageStatus | '')} value={status}><option value="">All</option><option value="published">Published</option><option value="draft">Draft</option><option value="archived">Archived</option></select>
        </div>
        <Link className="rounded-full bg-[var(--primary)] px-4 py-2 text-sm font-medium text-[var(--primary-foreground)]" href="/admin/pages/new">New page</Link>
      </Card>
      {pages.isLoading ? <LoadingState /> : null}
      {pages.isError ? <ApiErrorState message={pages.error.message} /> : null}
      {pages.data && pages.data.pages.length === 0 ? <EmptyState title="No pages" description="Create about, friends, projects, or other custom pages." /> : null}
      {pages.data?.pages.map((page) => (
        <Card className="grid gap-3 md:grid-cols-[1fr_auto] md:items-center" key={page.id}>
          <div><div className="flex items-center gap-2"><PageStatusBadge status={page.status} /><span className="text-xs text-[var(--muted-foreground)]">{page.template}</span></div><Link className="mt-2 block text-xl font-semibold hover:text-[var(--primary)]" href={`/admin/pages/${page.id}/edit`}>{page.title}</Link><p className="mt-1 text-sm text-[var(--muted-foreground)]">/{page.slug} · {page.viewCount} views</p></div>
          <div className="flex gap-2"><Link className="rounded-full border border-[var(--border)] px-3 py-1 text-xs" href={`/admin/pages/${page.id}/edit`}>Edit</Link><Button className="bg-red-600 px-3 py-1 text-xs text-white" disabled={remove.isPending} onClick={() => remove.mutate(page.id)} type="button">Delete</Button></div>
        </Card>
      ))}
    </div>
  );
}
