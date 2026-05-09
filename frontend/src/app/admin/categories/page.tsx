import { AdminCategories } from '@/components/features/taxonomy/admin-categories';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AdminCategoriesPage() {
  return <AdminShell><div className="mb-6"><h1 className="text-3xl font-semibold tracking-tight">Categories</h1><p className="mt-2 text-[var(--muted-foreground)]">Manage the category tree used by posts.</p></div><AdminCategories /></AdminShell>;
}
