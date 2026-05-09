import { PostList } from '@/components/features/post/post-list';
import { SiteContainer } from '@/components/layout/site-container';
import { SiteFooter } from '@/components/layout/site-footer';
import { SiteHeader } from '@/components/layout/site-header';
import { taxonomyApi } from '@/lib/api/taxonomy';

export const dynamic = 'force-dynamic';

export default async function TagPostsPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;
  const result = await taxonomyApi.tagPosts(slug, { perPage: 10 }).catch(() => ({ posts: [], meta: { page: 1, perPage: 10, total: 0, totalPages: 0 } }));
  return <><SiteHeader /><main className="py-12"><SiteContainer><div className="mb-8"><h1 className="text-4xl font-semibold tracking-tight">Tag: {slug}</h1><p className="mt-3 text-[var(--muted-foreground)]">Published posts with this tag.</p></div><PostList posts={result.posts} /></SiteContainer></main><SiteFooter /></>;
}
