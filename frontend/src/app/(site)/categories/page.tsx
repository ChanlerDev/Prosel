import { CategoryTree } from '@/components/features/taxonomy/category-tree';
import { SiteContainer } from '@/components/layout/site-container';
import { SiteFooter } from '@/components/layout/site-footer';
import { SiteHeader } from '@/components/layout/site-header';
import { taxonomyApi } from '@/lib/api/taxonomy';

export const revalidate = 60;

export default async function CategoriesPage() {
  const categories = await taxonomyApi.categories().catch(() => []);
  return <><SiteHeader /><main className="py-12"><SiteContainer><div className="mb-8"><h1 className="text-4xl font-semibold tracking-tight">Categories</h1><p className="mt-3 text-[var(--muted-foreground)]">Browse posts by category.</p></div><CategoryTree categories={categories} /></SiteContainer></main><SiteFooter /></>;
}
