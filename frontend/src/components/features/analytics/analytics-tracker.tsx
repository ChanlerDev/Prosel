'use client';

import { usePathname } from 'next/navigation';
import { useEffect } from 'react';

import { analyticsApi } from '@/lib/api/analytics';

type Props = {
  refType?: 'post' | 'note' | 'page';
  refId?: string;
};

export function AnalyticsTracker({ refType, refId }: Props) {
  const pathname = usePathname();

  useEffect(() => {
    if (!pathname || pathname.startsWith('/admin') || pathname.startsWith('/login')) return;
    analyticsApi.recordPageView({ path: pathname, refType, refId, referer: document.referrer }).catch(() => undefined);
  }, [pathname, refType, refId]);

  return null;
}
