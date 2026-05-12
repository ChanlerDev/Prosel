'use client';

import { Badge } from '@/components/ui/badge';
import { Card } from '@/components/ui/card';
import { usePublicAISummary } from '@/lib/ai/hooks';

export function AISummaryCard({ refId, language = 'zh' }: { refId: string; language?: string }) {
  const summary = usePublicAISummary('post', refId, language);

  if (summary.isLoading || summary.isError || !summary.data) return null;

  return (
    <Card className="mb-8 border-[var(--border)] bg-[var(--muted)]/30">
      <div className="mb-3 flex items-center justify-between gap-3">
        <h2 className="text-sm font-semibold uppercase tracking-[0.2em] text-[var(--muted-foreground)]">AI summary</h2>
        {summary.data.model ? <span className="text-xs text-[var(--muted-foreground)]">{summary.data.model}</span> : null}
      </div>
      <p className="leading-7 text-[var(--foreground)]">{summary.data.summary}</p>
      {summary.data.keywords.length > 0 ? (
        <div className="mt-4 flex flex-wrap gap-2">
          {summary.data.keywords.map((keyword) => <Badge key={keyword}>{keyword}</Badge>)}
        </div>
      ) : null}
    </Card>
  );
}
