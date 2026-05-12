'use client';

import { PageEditor } from '@/components/features/page/page-editor';
import { PageStatusBadge } from '@/components/features/page/page-status-badge';
import { ApiErrorState, LoadingState } from '@/components/features/system/states';
import { useAdminPage, useUpdatePage } from '@/lib/pages/hooks';

export function AdminPageEdit({ id }: { id: string }) {
  const page = useAdminPage(id);
  const update = useUpdatePage(id);

  if (page.isLoading) return <LoadingState />;
  if (page.isError || !page.data) return <ApiErrorState message="Unable to load page." />;

  return (
    <div className="grid gap-5">
      <div className="flex items-center gap-3"><PageStatusBadge status={page.data.status} /><span className="text-sm text-[var(--muted-foreground)]">{page.data.template}</span></div>
      <PageEditor error={update.isError ? update.error.message : undefined} isPending={update.isPending} onSubmit={(values) => update.mutate(values)} page={page.data} refId={page.data.id} />
    </div>
  );
}
