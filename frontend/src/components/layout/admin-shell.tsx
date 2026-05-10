'use client';

import Link from 'next/link';
import type { ReactNode } from 'react';

import { AuthGuard } from '@/components/features/auth/auth-guard';
import { Button } from '@/components/ui/button';
import { useLogout, useMe } from '@/lib/auth/hooks';

export function AdminShell({ children }: { children: ReactNode }) {
  return (
    <AuthGuard>
      <AdminFrame>{children}</AdminFrame>
    </AuthGuard>
  );
}

function AdminFrame({ children }: { children: ReactNode }) {
  const me = useMe();
  const logout = useLogout();

  return (
    <div className="min-h-screen bg-[var(--background)]">
      <aside className="fixed inset-y-0 left-0 hidden w-64 border-r border-[var(--border)] bg-[var(--card)] p-6 md:block">
        <Link className="text-lg font-semibold" href="/admin">
          Prosel Admin
        </Link>
        <nav className="mt-8 grid gap-3 text-sm text-[var(--muted-foreground)]">
          <Link href="/admin">Dashboard</Link>
          <Link href="/admin/posts">Posts</Link>
          <Link href="/admin/categories">Categories</Link>
          <Link href="/admin/tags">Tags</Link>
          <Link href="/admin/topics">Topics</Link>
          <Link href="/admin/comments">Comments</Link>
          <Link href="/admin/notes">Notes</Link>
          <Link href="/admin/activity">Activity</Link>
          <Link href="/admin/profile">Profile</Link>
          <Link href="/admin/security">Security</Link>
          <Link href="/admin/settings">Settings</Link>
        </nav>
      </aside>
      <main className="md:pl-64">
        <div className="border-b border-[var(--border)] bg-[var(--card)] px-6 py-3 md:hidden">
          <nav className="flex gap-4 overflow-x-auto text-sm text-[var(--muted-foreground)]">
            <Link href="/admin">Dashboard</Link>
            <Link href="/admin/posts">Posts</Link>
            <Link href="/admin/categories">Categories</Link>
            <Link href="/admin/tags">Tags</Link>
            <Link href="/admin/topics">Topics</Link>
            <Link href="/admin/comments">Comments</Link>
            <Link href="/admin/notes">Notes</Link>
          </nav>
        </div>
        <div className="flex items-center justify-between border-b border-[var(--border)] bg-[var(--card)] px-6 py-4">
          <p className="text-sm text-[var(--muted-foreground)]">Signed in as {me.data?.displayName ?? 'admin'}</p>
          <Button className="px-3 py-1 text-xs" disabled={logout.isPending} onClick={() => logout.mutate()} type="button">
            Sign out
          </Button>
        </div>
        <div className="p-6">{children}</div>
      </main>
    </div>
  );
}
