import { TopicCard } from '@/components/features/taxonomy/topic-card';
import { SiteContainer } from '@/components/layout/site-container';
import { SiteFooter } from '@/components/layout/site-footer';
import { SiteHeader } from '@/components/layout/site-header';
import { taxonomyApi } from '@/lib/api/taxonomy';

export const revalidate = 60;

export default async function TopicsPage() {
  const topics = await taxonomyApi.topics().catch(() => []);
  return <><SiteHeader /><main className="py-12"><SiteContainer><div className="mb-8"><h1 className="text-4xl font-semibold tracking-tight">Topics</h1><p className="mt-3 text-[var(--muted-foreground)]">Curated collections of posts.</p></div><div className="grid gap-5 md:grid-cols-2">{topics.map((topic) => <TopicCard key={topic.id} topic={topic} />)}</div></SiteContainer></main><SiteFooter /></>;
}
