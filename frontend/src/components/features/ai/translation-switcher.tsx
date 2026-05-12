'use client';

import { useState } from 'react';

import { PostContent } from '@/components/features/post/post-content';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { usePublicAITranslation } from '@/lib/ai/hooks';
import type { Post } from '@/types/post';

const languages = [
  { label: 'Original', value: '' },
  { label: 'English', value: 'en' },
  { label: 'Japanese', value: 'ja' },
];

export function TranslationSwitcher({ post }: { post: Post }) {
  const [language, setLanguage] = useState('');
  const [customLanguage, setCustomLanguage] = useState('');
  const translation = usePublicAITranslation('post', post.id, language);
  const displayPost: Post = language && translation.data ? { ...post, title: translation.data.title || post.title, excerpt: translation.data.summary || post.excerpt, contentMarkdown: translation.data.contentMarkdown } : post;

  return (
    <div className="grid gap-5">
      <Card className="flex flex-wrap items-center justify-between gap-3 py-3">
        <p className="text-sm text-[var(--muted-foreground)]">Read this article in another language when a translation is available.</p>
        <div className="flex flex-wrap gap-2">
          {languages.map((item) => (
            <Button className={`px-3 py-1 text-xs ${language === item.value ? '' : 'bg-[var(--secondary)] text-[var(--secondary-foreground)]'}`} key={item.value || 'original'} onClick={() => setLanguage(item.value)} type="button">
              {item.label}
            </Button>
          ))}
          <Input className="h-8 w-20 rounded-full px-3 py-1 text-xs" onChange={(event) => setCustomLanguage(event.target.value)} placeholder="lang" value={customLanguage} />
          <Button className="px-3 py-1 text-xs bg-[var(--secondary)] text-[var(--secondary-foreground)]" disabled={!customLanguage.trim()} onClick={() => setLanguage(customLanguage.trim().toLowerCase())} type="button">
            Go
          </Button>
        </div>
      </Card>
      {language && translation.isError ? <p className="text-sm text-[var(--muted-foreground)]">No {language} translation is available yet.</p> : null}
      {language && translation.isLoading ? <p className="text-sm text-[var(--muted-foreground)]">Loading translation...</p> : null}
      {language && translation.data ? (
        <Card className="border-[var(--border)] bg-[var(--muted)]/20">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-[var(--muted-foreground)]">Translated title</p>
          <h2 className="mt-2 text-2xl font-semibold">{translation.data.title || post.title}</h2>
          {translation.data.summary ? <p className="mt-3 text-sm leading-6 text-[var(--muted-foreground)]">{translation.data.summary}</p> : null}
        </Card>
      ) : null}
      <PostContent post={displayPost} />
    </div>
  );
}
