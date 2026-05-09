import Link from 'next/link';

import { Card } from '@/components/ui/card';
import type { Topic } from '@/types/taxonomy';

export function TopicCard({ topic }: { topic: Topic }) {
  return (
    <Card>
      <Link className="text-xl font-semibold hover:text-[var(--primary)]" href={`/topics/${topic.slug}`}>{topic.name}</Link>
      {topic.description ? <p className="mt-2 text-sm leading-6 text-[var(--muted-foreground)]">{topic.description}</p> : null}
    </Card>
  );
}
