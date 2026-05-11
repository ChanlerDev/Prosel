'use client';

import { useState } from 'react';

import { PageStatusBadge } from '@/components/features/page/page-status-badge';
import { ApiErrorState, LoadingState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { useAdminFriends, useCreateFriend, useDeleteFriend, useUpdateFriend } from '@/lib/pages/hooks';
import type { Friend, FriendStatus, FriendValues } from '@/types/page';

const emptyValues: FriendValues = { name: '', url: '', avatarUrl: '', description: '', status: 'active', sortOrder: 0 };

export function AdminFriends() {
  const [status, setStatus] = useState<FriendStatus | ''>('');
  const friends = useAdminFriends(status);
  const create = useCreateFriend();
  const update = useUpdateFriend();
  const remove = useDeleteFriend();
  const [editing, setEditing] = useState<Friend | null>(null);
  const [values, setValues] = useState<FriendValues>(emptyValues);

  function submit() {
    if (editing) update.mutate({ id: editing.id, body: values });
    else create.mutate(values);
    setEditing(null);
    setValues(emptyValues);
  }

  return (
    <div className="grid gap-5 lg:grid-cols-[360px_1fr]">
      <Card>
        <form className="grid gap-3" onSubmit={(event) => { event.preventDefault(); submit(); }}>
          <h2 className="font-semibold">{editing ? 'Edit friend' : 'New friend'}</h2>
          <Input onChange={(event) => setValues({ ...values, name: event.target.value })} placeholder="Name" required value={values.name} />
          <Input onChange={(event) => setValues({ ...values, url: event.target.value })} placeholder="https://example.com" required type="url" value={values.url} />
          <Input onChange={(event) => setValues({ ...values, avatarUrl: event.target.value })} placeholder="Avatar URL" value={values.avatarUrl} />
          <Textarea onChange={(event) => setValues({ ...values, description: event.target.value })} placeholder="Description" value={values.description} />
          <select className="rounded-xl border border-[var(--border)] bg-transparent px-3 py-2 text-sm" onChange={(event) => setValues({ ...values, status: event.target.value as FriendStatus })} value={values.status}><option value="active">Active</option><option value="pending">Pending</option><option value="hidden">Hidden</option></select>
          <Input onChange={(event) => setValues({ ...values, sortOrder: Number(event.target.value) })} placeholder="Sort order" type="number" value={values.sortOrder} />
          <Button disabled={create.isPending || update.isPending} type="submit">Save friend</Button>
          {editing ? <Button onClick={() => { setEditing(null); setValues(emptyValues); }} type="button">Cancel</Button> : null}
        </form>
      </Card>
      <div className="grid gap-3">
        <Card><select className="rounded-xl border border-[var(--border)] bg-transparent px-3 py-2 text-sm" onChange={(event) => setStatus(event.target.value as FriendStatus | '')} value={status}><option value="">All status</option><option value="active">Active</option><option value="pending">Pending</option><option value="hidden">Hidden</option></select></Card>
        {friends.isLoading ? <LoadingState /> : null}
        {friends.isError ? <ApiErrorState message={friends.error.message} /> : null}
        {friends.data?.map((friend) => (
          <Card className="flex flex-col justify-between gap-3 md:flex-row md:items-center" key={friend.id}>
            <div><div className="flex items-center gap-2"><PageStatusBadge status={friend.status} /><p className="font-semibold">{friend.name}</p></div><p className="mt-1 text-xs text-[var(--muted-foreground)]">{friend.url}</p>{friend.description ? <p className="mt-2 text-sm text-[var(--muted-foreground)]">{friend.description}</p> : null}</div>
            <div className="flex gap-2"><Button className="px-3 py-1 text-xs" onClick={() => { setEditing(friend); setValues({ name: friend.name, url: friend.url, avatarUrl: friend.avatarUrl ?? '', description: friend.description ?? '', status: friend.status, sortOrder: friend.sortOrder }); }} type="button">Edit</Button><Button className="bg-red-600 px-3 py-1 text-xs text-white" disabled={remove.isPending} onClick={() => remove.mutate(friend.id)} type="button">Delete</Button></div>
          </Card>
        ))}
      </div>
    </div>
  );
}
