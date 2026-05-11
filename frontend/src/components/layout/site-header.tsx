import Link from 'next/link';

import { SiteContainer } from '@/components/layout/site-container';

export function SiteHeader() {
  return (
    <header className="border-b border-[var(--border)]">
      <SiteContainer>
        <div className="flex h-16 items-center justify-between">
          <Link className="font-semibold tracking-tight" href="/">
            Prosel
          </Link>
          <nav className="flex gap-5 text-sm text-[var(--muted-foreground)]">
            <Link href="/posts">Posts</Link>
            <Link href="/categories">Categories</Link>
            <Link href="/topics">Topics</Link>
            <Link href="/search">Search</Link>
            <Link href="/login">Login</Link>
          </nav>
        </div>
      </SiteContainer>
    </header>
  );
}
