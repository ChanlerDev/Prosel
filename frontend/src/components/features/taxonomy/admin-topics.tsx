'use client';

import { useState } from 'react';

import { ApiErrorState, LoadingState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { useCreateTopic, useDeleteTopic, useTopics, useUpdateTopic } from '@/lib/taxonomy/hooks';
import type { Topic, TopicValues } from '@/types/taxonomy';

const emptyValues: TopicValues = { name: '', slug: '', description: '', coverImage: '', sortOrder: 0, items: [] };

export function AdminTopics() {
  const topics = useTopics();
  const create = useCreateTopic();
  const update = useUpdateTopic();
  const remove = useDeleteTopic();
  const [editing, setEditing] = useState<Topic | null>(null);
  const [values, setValues] = useState(emptyValues);

  return (
    <div className="grid gap-5 lg:grid-cols-[360px_1fr]">
      <Card>
        <form className="grid gap-3" onSubmit={(event) => { event.preventDefault(); if (editing) { update.mutate({ id: editing.id, body: values }); } else { create.mutate(values); } setEditing(null); setValues(emptyValues); }}>
          <h2 className="font-semibold">{editing ? 'Edit topic' : 'New topic'}</h2>
          <Input onChange={(event) => setValues({ ...values, name: event.target.value })} placeholder="Name" value={values.name} />
          <Input onChange={(event) => setValues({ ...values, slug: event.target.value })} placeholder="Slug" value={values.slug} />
          <Input onChange={(event) => setValues({ ...values, coverImage: event.target.value })} placeholder="Cover image" value={values.coverImage} />
          <Input onChange={(event) => setValues({ ...values, sortOrder: Number(event.target.value) })} placeholder="Sort order" type="number" value={values.sortOrder} />
          <Textarea onChange={(event) => setValues({ ...values, description: event.target.value })} placeholder="Description" value={values.description} />
          <Button disabled={create.isPending || update.isPending} type="submit">Save topic</Button>
        </form>
      </Card>
      <div className="grid gap-3">
        {topics.isLoading ? <LoadingState /> : null}
        {topics.isError ? <ApiErrorState message={topics.error.message} /> : null}
        {topics.data?.map((topic) => (
          <Card className="flex items-center justify-between gap-3" key={topic.id}>
            <div><p className="font-semibold">{topic.name}</p><p className="text-xs text-[var(--muted-foreground)]">/{topic.slug}</p></div>
            <div className="flex gap-2"><Button className="px-3 py-1 text-xs" onClick={() => { setEditing(topic); setValues({ name: topic.name, slug: topic.slug, description: topic.description ?? '', coverImage: topic.coverImage ?? '', sortOrder: topic.sortOrder, items: topic.items ?? [] }); }} type="button">Edit</Button><Button className="bg-red-600 px-3 py-1 text-xs text-white" onClick={() => remove.mutate(topic.id)} type="button">Delete</Button></div>
          </Card>
        ))}
      </div>
    </div>
  );
}
