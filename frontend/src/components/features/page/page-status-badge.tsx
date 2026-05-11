import type { FriendStatus, PageStatus } from '@/types/page';

const statusClass: Record<PageStatus | FriendStatus, string> = {
  draft: 'bg-zinc-100 text-zinc-700',
  published: 'bg-emerald-100 text-emerald-800',
  archived: 'bg-slate-100 text-slate-700',
  active: 'bg-emerald-100 text-emerald-800',
  pending: 'bg-amber-100 text-amber-900',
  hidden: 'bg-zinc-100 text-zinc-700',
};

export function PageStatusBadge({ status }: { status: PageStatus | FriendStatus }) {
  return <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${statusClass[status]}`}>{status}</span>;
}
