import Link from 'next/link';

import { PostList } from '@/components/features/post/post-list';
import { SiteContainer } from '@/components/layout/site-container';
import { SiteFooter } from '@/components/layout/site-footer';
import { SiteHeader } from '@/components/layout/site-header';
import { Badge } from '@/components/ui/badge';
import { Card } from '@/components/ui/card';
import { api } from '@/lib/api/client';
import { postsApi } from '@/lib/api/posts';

export const revalidate = 60;

export default async function HomePage() {
  const [settings, featured, latest] = await Promise.all([
    api.settings.public().catch(() => null),
    postsApi.list({ featured: true, perPage: 3 }).catch(() => ({ posts: [] })),
    postsApi.list({ perPage: 6 }).catch(() => ({ posts: [] })),
  ]);
  const siteName = settings?.site_name ?? 'Prosel';
  const description = settings?.site_description ?? 'A personal blog powered by Prosel';

  return (
    <>
      <SiteHeader />
      <main className="py-16">
        <SiteContainer>
          <div className="max-w-3xl">
            <Badge>Project setup</Badge>
            <h1 className="mt-6 text-5xl font-semibold tracking-tight">{siteName}</h1>
            <p className="mt-5 text-lg leading-8 text-[var(--muted-foreground)]">{description}</p>
          </div>
          <div className="mt-10 grid gap-4 md:grid-cols-3">
            <Card>
              <h2 className="font-semibold">Published posts</h2>
              <p className="mt-2 text-sm text-[var(--muted-foreground)]">Read the latest long-form writing and project updates.</p>
            </Card>
            <Card>
              <h2 className="font-semibold">Markdown first</h2>
              <p className="mt-2 text-sm text-[var(--muted-foreground)]">Posts are written in Markdown with SEO metadata and read counts.</p>
            </Card>
            <Card>
              <h2 className="font-semibold">Admin publishing</h2>
              <p className="mt-2 text-sm text-[var(--muted-foreground)]">Draft, edit, publish, and unpublish from the protected dashboard.</p>
            </Card>
          </div>
          {featured.posts.length > 0 ? (
            <section className="mt-14">
              <div className="mb-5 flex items-end justify-between">
                <h2 className="text-2xl font-semibold">Featured</h2>
                <Link className="text-sm text-[var(--primary)]" href="/posts">All posts</Link>
              </div>
              <PostList posts={featured.posts} />
            </section>
          ) : null}
          <section className="mt-14">
            <div className="mb-5 flex items-end justify-between">
              <h2 className="text-2xl font-semibold">Latest posts</h2>
              <Link className="text-sm text-[var(--primary)]" href="/posts">All posts</Link>
            </div>
            <PostList posts={latest.posts} />
          </section>
        </SiteContainer>
      </main>
      <SiteFooter />
    </>
  );
}
