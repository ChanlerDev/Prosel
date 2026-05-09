import Link from 'next/link';

import type { Tag } from '@/types/taxonomy';

export function TagChips({ tags }: { tags: Tag[] }) {
  if (tags.length === 0) return <p className="text-sm text-[var(--muted-foreground)]">No tags yet.</p>;
  return (
    <div className="flex flex-wrap gap-2">
      {tags.map((tag) => (
        <Link className="rounded-full border border-[var(--border)] px-3 py-1 text-sm" href={`/tags/${tag.slug}`} key={tag.id} style={tag.color ? { borderColor: tag.color } : undefined}>
          {tag.name}{tag.postCount !== undefined ? ` · ${tag.postCount}` : ''}
        </Link>
      ))}
    </div>
  );
}
