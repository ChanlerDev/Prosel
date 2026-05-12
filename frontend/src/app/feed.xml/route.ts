import { env } from '@/lib/env';

export const dynamic = 'force-dynamic';

export async function GET() {
  const response = await fetch(`${env.apiBaseUrl.replace(/\/api\/v1$/, '')}/feed.xml`, { cache: 'no-store' });
  return new Response(await response.text(), {
    status: response.status,
    headers: { 'Content-Type': 'application/rss+xml; charset=utf-8' },
  });
}
