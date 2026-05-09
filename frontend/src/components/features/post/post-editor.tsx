'use client';

import { useState } from 'react';

import { ApiErrorState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import type { Post, PostEditorValues } from '@/types/post';

const emptyValues: PostEditorValues = {
  title: '',
  slug: '',
  excerpt: '',
  contentMarkdown: '',
  coverImage: '',
  categoryId: '',
  featured: false,
  seoTitle: '',
  seoDescription: '',
};

export function PostEditor({ post, isPending, error, onSubmit }: { post?: Post; isPending: boolean; error?: string; onSubmit: (values: PostEditorValues) => void }) {
  const [values, setValues] = useState<PostEditorValues>(() => (post ? valuesFromPost(post) : emptyValues));

  return (
    <Card>
      <form
        className="grid gap-5"
        onSubmit={(event) => {
          event.preventDefault();
          onSubmit(values);
        }}
      >
        <div className="grid gap-4 md:grid-cols-2">
          <label className="grid gap-2 text-sm">
            Title
            <Input onChange={(event) => update('title', event.target.value)} value={values.title} />
          </label>
          <label className="grid gap-2 text-sm">
            Slug
            <Input onChange={(event) => update('slug', event.target.value)} placeholder="auto-generated if empty" value={values.slug} />
          </label>
        </div>
        <label className="grid gap-2 text-sm">
          Excerpt
          <Textarea onChange={(event) => update('excerpt', event.target.value)} rows={3} value={values.excerpt} />
        </label>
        <label className="grid gap-2 text-sm">
          Markdown
          <Textarea className="font-mono" onChange={(event) => update('contentMarkdown', event.target.value)} rows={16} value={values.contentMarkdown} />
        </label>
        <div className="grid gap-4 md:grid-cols-2">
          <label className="grid gap-2 text-sm">
            Cover image URL
            <Input onChange={(event) => update('coverImage', event.target.value)} value={values.coverImage} />
          </label>
          <label className="grid gap-2 text-sm">
            SEO title
            <Input onChange={(event) => update('seoTitle', event.target.value)} value={values.seoTitle} />
          </label>
        </div>
        <label className="grid gap-2 text-sm">
          SEO description
          <Textarea onChange={(event) => update('seoDescription', event.target.value)} rows={3} value={values.seoDescription} />
        </label>
        <label className="flex items-center gap-2 text-sm">
          <input checked={values.featured} onChange={(event) => update('featured', event.target.checked)} type="checkbox" />
          Featured post
        </label>
        {error ? <ApiErrorState message={error} /> : null}
        <Button className="w-fit" disabled={isPending} type="submit">
          {isPending ? 'Saving...' : 'Save post'}
        </Button>
      </form>
    </Card>
  );

  function update<Key extends keyof PostEditorValues>(key: Key, value: PostEditorValues[Key]) {
    setValues((current) => ({ ...current, [key]: value }));
  }
}

function valuesFromPost(post: Post): PostEditorValues {
  return {
    title: post.title,
    slug: post.slug,
    excerpt: post.excerpt ?? '',
    contentMarkdown: post.contentMarkdown ?? '',
    coverImage: post.coverImage ?? '',
    categoryId: post.categoryId ?? '',
    featured: post.featured,
    seoTitle: post.seoTitle ?? '',
    seoDescription: post.seoDescription ?? '',
  };
}
