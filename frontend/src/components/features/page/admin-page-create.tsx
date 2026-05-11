'use client';

import { PageEditor } from '@/components/features/page/page-editor';
import { useCreatePage } from '@/lib/pages/hooks';

export function AdminPageCreate() {
  const create = useCreatePage();
  return <PageEditor error={create.isError ? create.error.message : undefined} isPending={create.isPending} onSubmit={(values) => create.mutate(values)} />;
}
