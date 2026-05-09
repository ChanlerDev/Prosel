'use client';

import Link from 'next/link';
import { useState } from 'react';

import { ApiErrorState, EmptyState, LoadingState } from '@/components/features/system/states';
import { PostPublishButton } from '@/components/features/post/post-publish-button';
import { PostStatusBadge } from '@/components/features/post/post-status-badge';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { useAdminPosts, useDeletePost } from '@/lib/posts/hooks';
import type { PostStatus } from '@/types/post';

export function AdminPostsList() {
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState('');
  const posts = useAdminPosts({ search, status: status as PostStatus | '' });
  const deletePost = useDeletePost();

  return (
    <div className="grid gap-5">
      <Card className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div className="flex gap-3">
          <Input onChange={(event) => setSearch(event.target.value)} placeholder="Search posts" value={search} />
          <select className="rounded-xl border border-[var(--border)] bg-transparent px-3 py-2 text-sm" onChange={(event) => setStatus(event.target.value)} value={status}>
            <option value="">All</option>
            <option value="draft">Draft</option>
            <option value="published">Published</option>
            <option value="archived">Archived</option>
          </select>
        </div>
        <Link className="rounded-full bg-[var(--primary)] px-4 py-2 text-sm font-medium text-[var(--primary-foreground)]" href="/admin/posts/new">
          New post
        </Link>
      </Card>
      {posts.isLoading ? <LoadingState /> : null}
      {posts.isError ? <ApiErrorState message={posts.error.message} /> : null}
      {posts.data && posts.data.posts.length === 0 ? <EmptyState title="No posts" description="Create your first draft to start the publishing flow." /> : null}
      {posts.data?.posts.map((post) => (
        <Card className="grid gap-3 md:grid-cols-[1fr_auto] md:items-center" key={post.id}>
          <div>
            <div className="flex items-center gap-2">
              <PostStatusBadge status={post.status} />
              {post.featured ? <span className="text-xs text-[var(--muted-foreground)]">Featured</span> : null}
            </div>
            <Link className="mt-2 block text-xl font-semibold hover:text-[var(--primary)]" href={`/admin/posts/${post.id}/edit`}>
              {post.title}
            </Link>
            <p className="mt-1 text-sm text-[var(--muted-foreground)]">/{post.slug} · {post.viewCount} views</p>
          </div>
          <div className="flex gap-2">
            <PostPublishButton post={post} />
            <Link className="rounded-full border border-[var(--border)] px-3 py-1 text-xs" href={`/admin/posts/${post.id}/edit`}>
              Edit
            </Link>
            <Button className="bg-red-600 px-3 py-1 text-xs text-white" disabled={deletePost.isPending} onClick={() => deletePost.mutate(post.id)} type="button">
              Delete
            </Button>
          </div>
        </Card>
      ))}
    </div>
  );
}
