'use client';

import { useState } from 'react';

import { ApiErrorState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import type { Page, PageEditorValues, PageStatus, PageTemplate } from '@/types/page';

const emptyValues: PageEditorValues = { title: '', slug: '', subtitle: '', contentMarkdown: '', template: 'default', status: 'published', sortOrder: 0, seoTitle: '', seoDescription: '' };

export function PageEditor({ page, isPending, error, onSubmit }: { page?: Page; isPending: boolean; error?: string; onSubmit: (values: PageEditorValues) => void }) {
  const [values, setValues] = useState<PageEditorValues>(() => (page ? valuesFromPage(page) : emptyValues));

  return (
    <Card>
      <form className="grid gap-5" onSubmit={(event) => { event.preventDefault(); onSubmit(values); }}>
        <div className="grid gap-4 md:grid-cols-2">
          <label className="grid gap-2 text-sm">Title<Input onChange={(event) => update('title', event.target.value)} required value={values.title} /></label>
          <label className="grid gap-2 text-sm">Slug<Input onChange={(event) => update('slug', event.target.value)} placeholder="auto-generated if empty" value={values.slug} /></label>
        </div>
        <label className="grid gap-2 text-sm">Subtitle<Input onChange={(event) => update('subtitle', event.target.value)} value={values.subtitle} /></label>
        <label className="grid gap-2 text-sm">Markdown<Textarea className="font-mono" onChange={(event) => update('contentMarkdown', event.target.value)} required rows={14} value={values.contentMarkdown} /></label>
        <div className="grid gap-4 md:grid-cols-3">
          <label className="grid gap-2 text-sm">Template<select className="rounded-xl border border-[var(--border)] bg-transparent px-3 py-2 text-sm" onChange={(event) => update('template', event.target.value as PageTemplate)} value={values.template}><option value="default">Default</option><option value="about">About</option><option value="friends">Friends</option><option value="projects">Projects</option></select></label>
          <label className="grid gap-2 text-sm">Status<select className="rounded-xl border border-[var(--border)] bg-transparent px-3 py-2 text-sm" onChange={(event) => update('status', event.target.value as PageStatus)} value={values.status}><option value="published">Published</option><option value="draft">Draft</option><option value="archived">Archived</option></select></label>
          <label className="grid gap-2 text-sm">Sort order<Input onChange={(event) => update('sortOrder', Number(event.target.value))} type="number" value={values.sortOrder} /></label>
        </div>
        <div className="grid gap-4 md:grid-cols-2">
          <label className="grid gap-2 text-sm">SEO title<Input onChange={(event) => update('seoTitle', event.target.value)} value={values.seoTitle} /></label>
          <label className="grid gap-2 text-sm">SEO description<Textarea onChange={(event) => update('seoDescription', event.target.value)} rows={3} value={values.seoDescription} /></label>
        </div>
        {error ? <ApiErrorState message={error} /> : null}
        <Button className="w-fit" disabled={isPending} type="submit">{isPending ? 'Saving...' : 'Save page'}</Button>
      </form>
    </Card>
  );

  function update<Key extends keyof PageEditorValues>(key: Key, value: PageEditorValues[Key]) {
    setValues((current) => ({ ...current, [key]: value }));
  }
}

function valuesFromPage(page: Page): PageEditorValues {
  return { title: page.title, slug: page.slug, subtitle: page.subtitle ?? '', contentMarkdown: page.contentMarkdown ?? page.contentText, template: page.template, status: page.status, sortOrder: page.sortOrder, seoTitle: page.seoTitle ?? '', seoDescription: page.seoDescription ?? '' };
}
