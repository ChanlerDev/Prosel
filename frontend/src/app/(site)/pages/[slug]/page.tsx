import type { Metadata } from 'next';
import { notFound } from 'next/navigation';

import { PageContent } from '@/components/features/page/page-content';
import { SiteContainer } from '@/components/layout/site-container';
import { SiteFooter } from '@/components/layout/site-footer';
import { SiteHeader } from '@/components/layout/site-header';
import { Badge } from '@/components/ui/badge';
import { pagesApi } from '@/lib/api/pages';

export const revalidate = 60;

export async function generateMetadata({ params }: { params: Promise<{ slug: string }> }): Promise<Metadata> {
  const { slug } = await params;
  const page = await pagesApi.detail(slug).catch(() => null);
  if (!page) return {};
  return { title: page.seoTitle || page.title, description: page.seoDescription || page.subtitle || page.contentText.slice(0, 160) };
}

export default async function CustomPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;
  const page = await pagesApi.detail(slug).catch(() => null);
  if (!page) notFound();

  return (
    <>
      <SiteHeader />
      <main className="py-12">
        <SiteContainer>
          <article className="mx-auto max-w-3xl">
            <Badge>{page.template}</Badge>
            <h1 className="mt-5 text-5xl font-semibold tracking-tight">{page.title}</h1>
            {page.subtitle ? <p className="mt-5 text-lg leading-8 text-[var(--muted-foreground)]">{page.subtitle}</p> : null}
            <div className="mt-8"><PageContent page={page} /></div>
          </article>
        </SiteContainer>
      </main>
      <SiteFooter />
    </>
  );
}
