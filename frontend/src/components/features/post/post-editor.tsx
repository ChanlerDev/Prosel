'use client';

import { useRef, useState } from 'react';

import { MediaPicker } from '@/components/features/file/media-picker';
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
  tagIds: [],
  featured: false,
  seoTitle: '',
  seoDescription: '',
};

export function PostEditor({ post, refId, isPending, error, onSubmit }: { post?: Post; refId?: string; isPending: boolean; error?: string; onSubmit: (values: PostEditorValues) => void }) {
  const [values, setValues] = useState<PostEditorValues>(() => (post ? valuesFromPost(post) : emptyValues));
  const markdownRef = useRef<HTMLTextAreaElement>(null);

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
          <Textarea ref={markdownRef} className="font-mono" onChange={(event) => update('contentMarkdown', event.target.value)} rows={16} value={values.contentMarkdown} />
        </label>
        <MediaPicker refId={refId} refType="post" onSelect={(file) => insertMarkdownImage(file.originalName, file.publicUrl)} />
        <div className="grid gap-4 md:grid-cols-2">
          <label className="grid gap-2 text-sm">
            Category ID
            <Input onChange={(event) => update('categoryId', event.target.value)} value={values.categoryId} />
          </label>
          <label className="grid gap-2 text-sm">
            Tag IDs
            <Input onChange={(event) => update('tagIds', event.target.value.split(',').map((item) => item.trim()).filter(Boolean))} placeholder="tag-id-1, tag-id-2" value={values.tagIds.join(', ')} />
          </label>
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

  function insertMarkdownImage(alt: string, url: string) {
    const textarea = markdownRef.current;
    const markdown = `![${alt}](${url})`;
    const start = textarea?.selectionStart ?? values.contentMarkdown.length;
    const end = textarea?.selectionEnd ?? values.contentMarkdown.length;
    const next = values.contentMarkdown.slice(0, start) + markdown + values.contentMarkdown.slice(end);
    update('contentMarkdown', next);
    requestAnimationFrame(() => {
      textarea?.focus();
      textarea?.setSelectionRange(start + markdown.length, start + markdown.length);
    });
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
    tagIds: post.tagIds ?? [],
    featured: post.featured,
    seoTitle: post.seoTitle ?? '',
    seoDescription: post.seoDescription ?? '',
  };
}
