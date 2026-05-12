'use client';

import { useState } from 'react';

import { ApiErrorState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { useAIStatus, useGenerateSummary, useGenerateTranslation } from '@/lib/ai/hooks';
import type { Post } from '@/types/post';

export function AdminAITools({ post }: { post: Post }) {
  const status = useAIStatus();
  const summary = useGenerateSummary();
  const translation = useGenerateTranslation();
  const [summaryLanguage, setSummaryLanguage] = useState('zh');
  const [sourceLanguage, setSourceLanguage] = useState('zh');
  const [targetLanguage, setTargetLanguage] = useState('en');

  if (status.isLoading) return null;
  if (status.isError) return <ApiErrorState message={status.error.message} />;
  if (!status.data?.configured) {
    return (
      <Card className="border-amber-300 bg-amber-50 text-amber-950">
        <h2 className="font-semibold">AI provider is not configured</h2>
        <p className="mt-2 text-sm">Set AI_PROVIDER=openai and AI_API_KEY to enable summary and translation generation.</p>
      </Card>
    );
  }

  return (
    <Card className="grid gap-5">
      <div>
        <h2 className="text-lg font-semibold">AI tools</h2>
        <p className="mt-1 text-sm text-[var(--muted-foreground)]">Generate and save summaries or translations for the public article page.</p>
      </div>
      <div className="grid gap-4 md:grid-cols-2">
        <div className="grid gap-3 rounded-xl border border-[var(--border)] p-4">
          <label className="grid gap-2 text-sm">
            Summary language
            <Input onChange={(event) => setSummaryLanguage(event.target.value)} value={summaryLanguage} />
          </label>
          <Button disabled={summary.isPending} onClick={() => summary.mutate({ refType: 'post', refId: post.id, language: summaryLanguage })} type="button">
            {summary.isPending ? 'Generating...' : 'Generate summary'}
          </Button>
          {summary.isError ? <p className="text-sm text-red-600">{summary.error.message}</p> : null}
          {summary.data ? <p className="text-sm text-[var(--muted-foreground)]">Saved summary: {summary.data.summary}</p> : null}
        </div>
        <div className="grid gap-3 rounded-xl border border-[var(--border)] p-4">
          <div className="grid gap-3 sm:grid-cols-2">
            <label className="grid gap-2 text-sm">
              From
              <Input onChange={(event) => setSourceLanguage(event.target.value)} value={sourceLanguage} />
            </label>
            <label className="grid gap-2 text-sm">
              To
              <Input onChange={(event) => setTargetLanguage(event.target.value)} value={targetLanguage} />
            </label>
          </div>
          <Button disabled={translation.isPending} onClick={() => translation.mutate({ refType: 'post', refId: post.id, sourceLanguage, targetLanguage })} type="button">
            {translation.isPending ? 'Generating...' : 'Generate translation'}
          </Button>
          {translation.isError ? <p className="text-sm text-red-600">{translation.error.message}</p> : null}
          {translation.data ? <p className="text-sm text-[var(--muted-foreground)]">Saved {translation.data.targetLanguage} translation.</p> : null}
        </div>
      </div>
    </Card>
  );
}
