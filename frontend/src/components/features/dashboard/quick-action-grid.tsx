import Link from 'next/link';

import { Card } from '@/components/ui/card';

const actions = [
  { title: 'Write post', href: '/admin/posts/new', description: 'Create a Markdown draft' },
  { title: 'Manage posts', href: '/admin/posts', description: 'Edit, publish, and unpublish' },
  { title: 'Categories', href: '/admin/categories', description: 'Organize posts by tree' },
  { title: 'Tags', href: '/admin/tags', description: 'Manage tag chips' },
  { title: 'Topics', href: '/admin/topics', description: 'Curate collections' },
];

export function QuickActionGrid() {
  return (
    <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-3">
      {actions.map((action) => (
        <Link href={action.href} key={action.href}>
          <Card className="h-full transition hover:border-[var(--primary)]">
            <h3 className="font-semibold">{action.title}</h3>
            <p className="mt-2 text-sm text-[var(--muted-foreground)]">{action.description}</p>
          </Card>
        </Link>
      ))}
    </div>
  );
}
