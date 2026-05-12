'use client';

import { useRef, useState } from 'react';

import { MediaPicker } from '@/components/features/file/media-picker';
import { ApiErrorState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import type { Note, NoteEditorValues, NoteStatus } from '@/types/note';

const emptyValues: NoteEditorValues = { title: '', slug: '', contentMarkdown: '', mood: '', weather: '', location: '', status: 'published' };

export function NoteEditor({ note, refId, isPending, error, onSubmit }: { note?: Note; refId?: string; isPending: boolean; error?: string; onSubmit: (values: NoteEditorValues) => void }) {
  const [values, setValues] = useState<NoteEditorValues>(() => (note ? valuesFromNote(note) : emptyValues));
  const markdownRef = useRef<HTMLTextAreaElement>(null);

  return (
    <Card>
      <form className="grid gap-5" onSubmit={(event) => { event.preventDefault(); onSubmit(values); }}>
        <div className="grid gap-4 md:grid-cols-2">
          <label className="grid gap-2 text-sm">
            Title
            <Input onChange={(event) => update('title', event.target.value)} placeholder="Optional title" value={values.title} />
          </label>
          <label className="grid gap-2 text-sm">
            Slug
            <Input onChange={(event) => update('slug', event.target.value)} placeholder="auto-generated if empty" value={values.slug} />
          </label>
        </div>
        <label className="grid gap-2 text-sm">
          Markdown
          <Textarea ref={markdownRef} className="font-mono" onChange={(event) => update('contentMarkdown', event.target.value)} required rows={12} value={values.contentMarkdown} />
        </label>
        <MediaPicker refId={refId} refType="note" onSelect={(file) => insertMarkdownImage(file.originalName, file.publicUrl)} />
        <div className="grid gap-4 md:grid-cols-4">
          <label className="grid gap-2 text-sm">
            Mood
            <Input onChange={(event) => update('mood', event.target.value)} placeholder="calm" value={values.mood} />
          </label>
          <label className="grid gap-2 text-sm">
            Weather
            <Input onChange={(event) => update('weather', event.target.value)} placeholder="sunny" value={values.weather} />
          </label>
          <label className="grid gap-2 text-sm">
            Location
            <Input onChange={(event) => update('location', event.target.value)} placeholder="home" value={values.location} />
          </label>
          <label className="grid gap-2 text-sm">
            Status
            <select className="rounded-xl border border-[var(--border)] bg-transparent px-3 py-2 text-sm" onChange={(event) => update('status', event.target.value as NoteStatus)} value={values.status}>
              <option value="published">Published</option>
              <option value="draft">Draft</option>
              <option value="private">Private</option>
              <option value="archived">Archived</option>
            </select>
          </label>
        </div>
        {error ? <ApiErrorState message={error} /> : null}
        <Button className="w-fit" disabled={isPending} type="submit">{isPending ? 'Saving...' : 'Save note'}</Button>
      </form>
    </Card>
  );

  function update<Key extends keyof NoteEditorValues>(key: Key, value: NoteEditorValues[Key]) {
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

function valuesFromNote(note: Note): NoteEditorValues {
  return { title: note.title ?? '', slug: note.slug, contentMarkdown: note.contentMarkdown ?? note.contentText, mood: note.mood ?? '', weather: note.weather ?? '', location: note.location ?? '', status: note.status };
}
