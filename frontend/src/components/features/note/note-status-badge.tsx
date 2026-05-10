import type { NoteStatus } from '@/types/note';

const statusClass: Record<NoteStatus, string> = {
  draft: 'bg-zinc-100 text-zinc-700',
  published: 'bg-emerald-100 text-emerald-800',
  private: 'bg-amber-100 text-amber-900',
  archived: 'bg-slate-100 text-slate-700',
};

export function NoteStatusBadge({ status }: { status: NoteStatus }) {
  return <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${statusClass[status]}`}>{status}</span>;
}
