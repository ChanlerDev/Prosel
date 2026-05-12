import Link from 'next/link';

import { Card } from '@/components/ui/card';
import { subscribeApi } from '@/lib/api/subscribe';

export default async function UnsubscribePage({ searchParams }: { searchParams: Promise<{ token?: string }> }) {
  const { token } = await searchParams;
  let message = 'You have been unsubscribed.';
  let ok = true;
  try {
    await subscribeApi.unsubscribe(token ?? '');
  } catch (error) {
    ok = false;
    message = error instanceof Error ? error.message : 'Unsubscribe failed.';
  }

  return (
    <main className="mx-auto grid min-h-screen max-w-xl place-items-center px-6">
      <Card>
        <h1 className="text-2xl font-semibold">{ok ? 'Unsubscribed' : 'Unsubscribe failed'}</h1>
        <p className="mt-3 text-[var(--muted-foreground)]">{message}</p>
        <Link className="mt-6 inline-flex rounded-full bg-[var(--primary)] px-4 py-2 text-sm font-medium text-[var(--primary-foreground)]" href="/">
          Back home
        </Link>
      </Card>
    </main>
  );
}
