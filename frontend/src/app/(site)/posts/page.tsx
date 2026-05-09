import { PostList } from '@/components/features/post/post-list';
import { SiteContainer } from '@/components/layout/site-container';
import { SiteFooter } from '@/components/layout/site-footer';
import { SiteHeader } from '@/components/layout/site-header';
import { postsApi } from '@/lib/api/posts';

export const dynamic = 'force-dynamic';

export default async function PostsPage({ searchParams }: { searchParams: Promise<{ page?: string }> }) {
  const params = await searchParams;
  const page = Number(params.page ?? '1');
  const result = await postsApi.list({ page: Number.isFinite(page) ? page : 1, perPage: 10 }).catch(() => ({ posts: [], meta: { page: 1, perPage: 10, total: 0, totalPages: 0 } }));

  return (
    <>
      <SiteHeader />
      <main className="py-12">
        <SiteContainer>
          <div className="mb-8">
            <h1 className="text-4xl font-semibold tracking-tight">Posts</h1>
            <p className="mt-3 text-[var(--muted-foreground)]">Published essays and notes from the blog.</p>
          </div>
          <PostList posts={result.posts} />
        </SiteContainer>
      </main>
      <SiteFooter />
    </>
  );
}
