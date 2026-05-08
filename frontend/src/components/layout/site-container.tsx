import type { ReactNode } from 'react';

export function SiteContainer({ children }: { children: ReactNode }) {
  return <div className="mx-auto w-full max-w-5xl px-6">{children}</div>;
}
