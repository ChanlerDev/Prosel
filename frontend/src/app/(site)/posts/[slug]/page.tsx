import type { Metadata } from 'next';
import { notFound } from 'next/navigation';

import { AISummaryCard } from '@/components/features/ai/ai-summary-card';
import { TranslationSwitcher } from '@/components/features/ai/translation-switcher';
import { CommentSection } from '@/components/features/comment/comment-section';
import { SiteContainer } from '@/components/layout/site-container';
import { SiteFooter } from '@/components/layout/site-footer';
import { SiteHeader } from '@/components/layout/site-header';
import { Badge } from '@/components/ui/badge';
import { postsApi } from '@/lib/api/posts';

export const revalidate = 60;

export async function generateMetadata({ params }: { params: Promise<{ slug: string }> }): Promise<Metadata> {
  const { slug } = await params;
  const post = await postsApi.detail(slug).catch(() => null);
  if (!post) return {};
  return {
    title: post.seoTitle || post.title,
    description: post.seoDescription || post.excerpt,
  };
}

export default async function PostDetailPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;
  const post = await postsApi.detail(slug).catch(() => null);
  if (!post) notFound();

  return (
    <>
      <SiteHeader />
      <main className="py-12">
        <SiteContainer>
          <article className="mx-auto max-w-3xl">
            <div className="mb-8">
              <Badge>{post.viewCount} views</Badge>
              <h1 className="mt-5 text-5xl font-semibold tracking-tight">{post.title}</h1>
              {post.excerpt ? <p className="mt-5 text-lg leading-8 text-[var(--muted-foreground)]">{post.excerpt}</p> : null}
              <p className="mt-4 text-sm text-[var(--muted-foreground)]">Published {formatDate(post.publishedAt ?? post.createdAt)}</p>
            </div>
            <AISummaryCard refId={post.id} />
            <TranslationSwitcher post={post} />
          </article>
          <CommentSection refId={post.id} refType="post" />
        </SiteContainer>
      </main>
      <SiteFooter />
    </>
  );
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat('en', { dateStyle: 'medium' }).format(new Date(value));
}
