import Link from 'next/link';

import { Card } from '@/components/ui/card';
import type { PostSummary } from '@/types/dashboard';

export function RecentPostTable({ posts }: { posts: PostSummary[] }) {
  if (posts.length === 0) {
    return <Card><p className="text-sm text-[var(--muted-foreground)]">No posts yet. Start with a new draft.</p></Card>;
  }
  return (
    <Card className="overflow-hidden p-0">
      <div className="border-b border-[var(--border)] p-4 font-semibold">Recent posts</div>
      <div className="divide-y divide-[var(--border)]">
        {posts.map((post) => (
          <div className="grid gap-2 p-4 md:grid-cols-[1fr_auto] md:items-center" key={post.id}>
            <div>
              <Link className="font-medium hover:text-[var(--primary)]" href={`/admin/posts/${post.id}/edit`}>{post.title}</Link>
              <p className="text-xs text-[var(--muted-foreground)]">{post.status} · {post.viewCount} views · updated {formatDate(post.updatedAt)}</p>
            </div>
            {post.status === 'published' ? <Link className="text-sm text-[var(--primary)]" href={`/posts/${post.slug}`}>View</Link> : null}
          </div>
        ))}
      </div>
    </Card>
  );
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat('en', { dateStyle: 'medium' }).format(new Date(value));
}
