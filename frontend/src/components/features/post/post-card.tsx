import Link from 'next/link';

import { Card } from '@/components/ui/card';
import type { Post } from '@/types/post';

export function PostCard({ post }: { post: Post }) {
  return (
    <Card className="grid gap-3">
      <div className="flex items-center gap-3 text-xs text-[var(--muted-foreground)]">
        {post.featured ? <span>Featured</span> : null}
        <time dateTime={post.publishedAt ?? post.createdAt}>{formatDate(post.publishedAt ?? post.createdAt)}</time>
        <span>{post.viewCount} views</span>
      </div>
      <Link className="text-2xl font-semibold tracking-tight hover:text-[var(--primary)]" href={`/posts/${post.slug}`}>
        {post.title}
      </Link>
      {post.excerpt ? <p className="text-sm leading-6 text-[var(--muted-foreground)]">{post.excerpt}</p> : null}
    </Card>
  );
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat('en', { dateStyle: 'medium' }).format(new Date(value));
}
