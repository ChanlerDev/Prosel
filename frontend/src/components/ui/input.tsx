import type { InputHTMLAttributes } from 'react';

import { cn } from '@/lib/utils';

export function Input({ className, ...props }: InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input
      className={cn(
        'w-full rounded-xl border border-[var(--border)] bg-transparent px-3 py-2 text-sm outline-none transition focus:border-[var(--primary)]',
        className,
      )}
      {...props}
    />
  );
}
