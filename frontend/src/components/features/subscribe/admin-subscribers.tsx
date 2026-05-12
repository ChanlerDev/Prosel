'use client';

import { useState } from 'react';

import { SubscriberStatusBadge } from '@/components/features/subscribe/subscriber-status-badge';
import { ApiErrorState, EmptyState, LoadingState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { useAdminPosts } from '@/lib/posts/hooks';
import { useAdminSubscribers, useNotifyPostSubscribers } from '@/lib/subscribe/hooks';
import type { SubscriberStatus } from '@/types/subscribe';

export function AdminSubscribers() {
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState<SubscriberStatus | ''>('');
  const [postId, setPostId] = useState('');
  const subscribers = useAdminSubscribers({ search, status });
  const posts = useAdminPosts({ status: 'published', perPage: 50 });
  const notify = useNotifyPostSubscribers();

  return (
    <div className="grid gap-5">
      <Card className="grid gap-3 md:grid-cols-[1fr_auto] md:items-end">
        <div className="grid gap-3 md:grid-cols-[1fr_auto]">
          <Input onChange={(event) => setSearch(event.target.value)} placeholder="Search subscribers" value={search} />
          <select className="rounded-xl border border-[var(--border)] bg-transparent px-3 py-2 text-sm" onChange={(event) => setStatus(event.target.value as SubscriberStatus | '')} value={status}>
            <option value="">All status</option>
            <option value="pending">Pending</option>
            <option value="active">Active</option>
            <option value="unsubscribed">Unsubscribed</option>
            <option value="bounced">Bounced</option>
          </select>
        </div>
        <div className="flex gap-2">
          <select className="min-w-56 rounded-xl border border-[var(--border)] bg-transparent px-3 py-2 text-sm" onChange={(event) => setPostId(event.target.value)} value={postId}>
            <option value="">Select published post</option>
            {posts.data?.posts.map((post) => (
              <option key={post.id} value={post.id}>{post.title}</option>
            ))}
          </select>
          <Button disabled={notify.isPending || !postId} onClick={() => notify.mutate(postId)} type="button">{notify.isPending ? 'Sending...' : 'Notify'}</Button>
        </div>
      </Card>
      {notify.isSuccess ? <Card className="border-emerald-200 bg-emerald-50 text-emerald-900">Notification sent to active subscribers.</Card> : null}
      {notify.isError ? <ApiErrorState message={notify.error.message} /> : null}
      {subscribers.isLoading ? <LoadingState /> : null}
      {subscribers.isError ? <ApiErrorState message={subscribers.error.message} /> : null}
      {subscribers.data && subscribers.data.subscribers.length === 0 ? <EmptyState title="No subscribers" description="Verified email subscribers will appear here." /> : null}
      {subscribers.data?.subscribers.map((subscriber) => (
        <Card className="grid gap-3 md:grid-cols-[1fr_auto] md:items-center" key={subscriber.id}>
          <div>
            <div className="flex flex-wrap items-center gap-2">
              <SubscriberStatusBadge status={subscriber.status} />
              <p className="font-semibold">{subscriber.email}</p>
              {subscriber.name ? <p className="text-sm text-[var(--muted-foreground)]">{subscriber.name}</p> : null}
            </div>
            <p className="mt-2 text-sm text-[var(--muted-foreground)]">Joined {formatDate(subscriber.createdAt)}{subscriber.verifiedAt ? ` · Verified ${formatDate(subscriber.verifiedAt)}` : ''}</p>
          </div>
        </Card>
      ))}
      {subscribers.data ? <p className="text-sm text-[var(--muted-foreground)]">Page {subscribers.data.meta.page} of {subscribers.data.meta.totalPages || 1} · {subscribers.data.meta.total} subscribers</p> : null}
    </div>
  );
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat('en', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value));
}
