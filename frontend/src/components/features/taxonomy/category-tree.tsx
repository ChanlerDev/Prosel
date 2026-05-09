import Link from 'next/link';

import { Card } from '@/components/ui/card';
import type { CategoryNode } from '@/types/taxonomy';

export function CategoryTree({ categories }: { categories: CategoryNode[] }) {
  if (categories.length === 0) return <Card>No categories yet.</Card>;
  return <div className="grid gap-3">{categories.map((category) => <CategoryBranch category={category} key={category.id} />)}</div>;
}

function CategoryBranch({ category }: { category: CategoryNode }) {
  return (
    <Card className="p-4">
      <Link className="font-semibold hover:text-[var(--primary)]" href={`/categories/${category.slug}`}>{category.name}</Link>
      <span className="ml-2 text-xs text-[var(--muted-foreground)]">{category.postCount} posts</span>
      {category.description ? <p className="mt-1 text-sm text-[var(--muted-foreground)]">{category.description}</p> : null}
      {category.children.length > 0 ? <div className="mt-3 grid gap-2 pl-4">{category.children.map((child) => <CategoryBranch category={child} key={child.id} />)}</div> : null}
    </Card>
  );
}
