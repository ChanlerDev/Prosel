import Link from 'next/link';
import type { ReactNode } from 'react';

export function AdminShell({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen bg-[var(--background)]">
      <aside className="fixed inset-y-0 left-0 hidden w-64 border-r border-[var(--border)] bg-[var(--card)] p-6 md:block">
        <Link className="text-lg font-semibold" href="/admin">
          Prosel Admin
        </Link>
        <nav className="mt-8 grid gap-3 text-sm text-[var(--muted-foreground)]">
          <Link href="/admin">Dashboard</Link>
          <Link href="/admin/posts">Posts</Link>
          <Link href="/admin/settings">Settings</Link>
        </nav>
      </aside>
      <main className="md:pl-64">
        <div className="border-b border-[var(--border)] bg-[var(--card)] px-6 py-4">
          <p className="text-sm text-[var(--muted-foreground)]">Admin workspace</p>
        </div>
        <div className="p-6">{children}</div>
      </main>
    </div>
  );
}
