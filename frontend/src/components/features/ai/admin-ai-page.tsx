'use client';

import { ApiErrorState, LoadingState } from '@/components/features/system/states';
import { Card } from '@/components/ui/card';
import { useAIStatus } from '@/lib/ai/hooks';

export function AdminAIPage() {
  const status = useAIStatus();

  if (status.isLoading) return <LoadingState />;
  if (status.isError) return <ApiErrorState message={status.error.message} />;

  return (
    <div className="grid gap-5">
      <Card>
        <h2 className="text-lg font-semibold">Provider status</h2>
        <p className="mt-2 text-sm text-[var(--muted-foreground)]">
          {status.data?.configured ? 'AI generation is enabled. Use each post edit page to generate summaries and translations.' : 'AI generation is disabled until AI_PROVIDER and AI_API_KEY are configured.'}
        </p>
      </Card>
      <Card>
        <h2 className="text-lg font-semibold">Configuration</h2>
        <ul className="mt-3 list-disc space-y-2 pl-5 text-sm text-[var(--muted-foreground)]">
          <li>AI_PROVIDER=openai</li>
          <li>AI_API_KEY=your provider key</li>
          <li>Optional: AI_BASE_URL, AI_MODEL, AI_TIMEOUT_SECONDS</li>
        </ul>
      </Card>
    </div>
  );
}
