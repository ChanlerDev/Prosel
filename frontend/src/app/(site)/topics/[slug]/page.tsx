import Link from 'next/link';

import { Card } from '@/components/ui/card';
import { SiteContainer } from '@/components/layout/site-container';
import { SiteFooter } from '@/components/layout/site-footer';
import { SiteHeader } from '@/components/layout/site-header';
import { taxonomyApi } from '@/lib/api/taxonomy';

export const dynamic = 'force-dynamic';

export default async function TopicDetailPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;
  const topic = await taxonomyApi.topic(slug).catch(() => null);
  return <><SiteHeader /><main className="py-12"><SiteContainer><div className="mb-8"><h1 className="text-4xl font-semibold tracking-tight">{topic?.name ?? slug}</h1>{topic?.description ? <p className="mt-3 text-[var(--muted-foreground)]">{topic.description}</p> : null}</div><div className="grid gap-3">{topic?.items?.map((item) => <Card key={`${item.refType}-${item.refId}`}><Link className="font-semibold hover:text-[var(--primary)]" href={item.slug ? `/posts/${item.slug}` : '#'}>{item.title || item.refId}</Link></Card>)}</div></SiteContainer></main><SiteFooter /></>;
}
