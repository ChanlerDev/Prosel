import { forwardRef, type TextareaHTMLAttributes } from 'react';

import { cn } from '@/lib/utils';

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaHTMLAttributes<HTMLTextAreaElement>>(function Textarea({ className, ...props }, ref) {
  return (
    <textarea
      ref={ref}
      className={cn(
        'w-full rounded-xl border border-[var(--border)] bg-transparent px-3 py-2 text-sm outline-none transition focus:border-[var(--primary)]',
        className,
      )}
      {...props}
    />
  );
});
