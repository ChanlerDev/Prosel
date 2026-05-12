import { AdminAIPage } from '@/components/features/ai/admin-ai-page';
import { AdminShell } from '@/components/layout/admin-shell';

export default function AIPage() {
  return (
    <AdminShell>
      <div className="mb-6">
        <h1 className="text-3xl font-semibold tracking-tight">AI</h1>
        <p className="mt-2 text-[var(--muted-foreground)]">Configure optional AI generation and use post edit pages to trigger summaries or translations.</p>
      </div>
      <AdminAIPage />
    </AdminShell>
  );
}
