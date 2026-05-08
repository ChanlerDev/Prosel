import type { HTMLAttributes } from 'react';

import { cn } from '@/lib/utils';

export function Badge({ className, ...props }: HTMLAttributes<HTMLSpanElement>) {
  return <span className={cn('inline-flex rounded-full bg-[var(--muted)] px-3 py-1 text-xs text-[var(--muted-foreground)]', className)} {...props} />;
}
