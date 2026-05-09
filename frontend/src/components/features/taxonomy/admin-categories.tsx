'use client';

import { useState } from 'react';

import { ApiErrorState, LoadingState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { useCategories, useCreateCategory, useDeleteCategory, useUpdateCategory } from '@/lib/taxonomy/hooks';
import type { CategoryNode, CategoryValues } from '@/types/taxonomy';

const emptyValues: CategoryValues = { parentId: '', name: '', slug: '', description: '', sortOrder: 0 };

export function AdminCategories() {
  const categories = useCategories();
  const create = useCreateCategory();
  const update = useUpdateCategory();
  const remove = useDeleteCategory();
  const [editing, setEditing] = useState<CategoryNode | null>(null);
  const [values, setValues] = useState(emptyValues);

  return (
    <div className="grid gap-5 lg:grid-cols-[360px_1fr]">
      <Card>
        <form className="grid gap-3" onSubmit={(event) => { event.preventDefault(); if (editing) { update.mutate({ id: editing.id, body: values }); } else { create.mutate(values); } setEditing(null); setValues(emptyValues); }}>
          <h2 className="font-semibold">{editing ? 'Edit category' : 'New category'}</h2>
          <Input onChange={(event) => setValues({ ...values, name: event.target.value })} placeholder="Name" value={values.name} />
          <Input onChange={(event) => setValues({ ...values, slug: event.target.value })} placeholder="Slug" value={values.slug} />
          <Input onChange={(event) => setValues({ ...values, parentId: event.target.value })} placeholder="Parent ID" value={values.parentId} />
          <Input onChange={(event) => setValues({ ...values, sortOrder: Number(event.target.value) })} placeholder="Sort order" type="number" value={values.sortOrder} />
          <Textarea onChange={(event) => setValues({ ...values, description: event.target.value })} placeholder="Description" value={values.description} />
          <Button disabled={create.isPending || update.isPending} type="submit">Save category</Button>
        </form>
      </Card>
      <div className="grid gap-3">
        {categories.isLoading ? <LoadingState /> : null}
        {categories.isError ? <ApiErrorState message={categories.error.message} /> : null}
        {categories.data?.map((category) => <CategoryAdminNode category={category} key={category.id} onDelete={(id) => remove.mutate(id)} onEdit={(node) => { setEditing(node); setValues({ parentId: node.parentId ?? '', name: node.name, slug: node.slug, description: node.description ?? '', sortOrder: node.sortOrder }); }} />)}
      </div>
    </div>
  );
}

function CategoryAdminNode({ category, onEdit, onDelete }: { category: CategoryNode; onEdit: (category: CategoryNode) => void; onDelete: (id: string) => void }) {
  return (
    <Card className="p-4">
      <div className="flex items-center justify-between gap-3">
        <div><p className="font-semibold">{category.name}</p><p className="text-xs text-[var(--muted-foreground)]">/{category.slug} · {category.id}</p></div>
        <div className="flex gap-2"><Button className="px-3 py-1 text-xs" onClick={() => onEdit(category)} type="button">Edit</Button><Button className="bg-red-600 px-3 py-1 text-xs text-white" onClick={() => onDelete(category.id)} type="button">Delete</Button></div>
      </div>
      {category.children.length > 0 ? <div className="mt-3 grid gap-2 pl-4">{category.children.map((child) => <CategoryAdminNode category={child} key={child.id} onDelete={onDelete} onEdit={onEdit} />)}</div> : null}
    </Card>
  );
}
