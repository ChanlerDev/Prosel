import { AdminTopics } from '@/components/features/taxonomy/admin-topics';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AdminTopicsPage() {
  return <AdminShell><div className="mb-6"><h1 className="text-3xl font-semibold tracking-tight">Topics</h1><p className="mt-2 text-[var(--muted-foreground)]">Manage curated topic collections.</p></div><AdminTopics /></AdminShell>;
}
