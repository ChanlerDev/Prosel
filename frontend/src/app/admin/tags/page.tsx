import { AdminTags } from '@/components/features/taxonomy/admin-tags';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AdminTagsPage() {
  return <AdminShell><div className="mb-6"><h1 className="text-3xl font-semibold tracking-tight">Tags</h1><p className="mt-2 text-[var(--muted-foreground)]">Manage tags for post discovery.</p></div><AdminTags /></AdminShell>;
}
