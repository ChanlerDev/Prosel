'use client';

import { useState } from 'react';

import { ApiErrorState, LoadingState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { useCreateTag, useDeleteTag, useTags, useUpdateTag } from '@/lib/taxonomy/hooks';
import type { Tag, TagValues } from '@/types/taxonomy';

const emptyValues: TagValues = { name: '', slug: '', color: '', description: '' };

export function AdminTags() {
  const tags = useTags();
  const create = useCreateTag();
  const update = useUpdateTag();
  const remove = useDeleteTag();
  const [editing, setEditing] = useState<Tag | null>(null);
  const [values, setValues] = useState(emptyValues);

  return (
    <div className="grid gap-5 lg:grid-cols-[360px_1fr]">
      <Card>
        <form className="grid gap-3" onSubmit={(event) => { event.preventDefault(); if (editing) { update.mutate({ id: editing.id, body: values }); } else { create.mutate(values); } setEditing(null); setValues(emptyValues); }}>
          <h2 className="font-semibold">{editing ? 'Edit tag' : 'New tag'}</h2>
          <Input onChange={(event) => setValues({ ...values, name: event.target.value })} placeholder="Name" value={values.name} />
          <Input onChange={(event) => setValues({ ...values, slug: event.target.value })} placeholder="Slug" value={values.slug} />
          <Input onChange={(event) => setValues({ ...values, color: event.target.value })} placeholder="#8a4f2a" value={values.color} />
          <Textarea onChange={(event) => setValues({ ...values, description: event.target.value })} placeholder="Description" value={values.description} />
          <Button disabled={create.isPending || update.isPending} type="submit">Save tag</Button>
        </form>
      </Card>
      <div className="grid gap-3">
        {tags.isLoading ? <LoadingState /> : null}
        {tags.isError ? <ApiErrorState message={tags.error.message} /> : null}
        {tags.data?.map((tag) => (
          <Card className="flex items-center justify-between gap-3" key={tag.id}>
            <div><p className="font-semibold">{tag.name}</p><p className="text-xs text-[var(--muted-foreground)]">/{tag.slug} · {tag.postCount ?? 0} posts</p></div>
            <div className="flex gap-2"><Button className="px-3 py-1 text-xs" onClick={() => { setEditing(tag); setValues({ name: tag.name, slug: tag.slug, color: tag.color ?? '', description: tag.description ?? '' }); }} type="button">Edit</Button><Button className="bg-red-600 px-3 py-1 text-xs text-white" onClick={() => remove.mutate(tag.id)} type="button">Delete</Button></div>
          </Card>
        ))}
      </div>
    </div>
  );
}
