'use client';

import { useState } from 'react';

import { ApiErrorState, EmptyState, LoadingState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { FileUploader } from '@/components/features/file/file-uploader';
import { useAdminFiles, useDeleteFile } from '@/lib/files/hooks';
import type { FileAsset, FileStatus } from '@/types/file';

export function MediaLibrary({ onSelect, compact = false }: { onSelect?: (file: FileAsset) => void; compact?: boolean }) {
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState<FileStatus | ''>('');
  const files = useAdminFiles({ search, status, type: 'image/' });
  const deleteFile = useDeleteFile();

  return (
    <div className="grid gap-5">
      <Card className="grid gap-4">
        <FileUploader onUploaded={onSelect} />
        <div className="flex flex-col gap-3 md:flex-row">
          <Input onChange={(event) => setSearch(event.target.value)} placeholder="Search files" value={search} />
          <select className="rounded-xl border border-[var(--border)] bg-transparent px-3 py-2 text-sm" onChange={(event) => setStatus(event.target.value as FileStatus | '')} value={status}>
            <option value="">Active files</option>
            <option value="orphan">Orphan</option>
            <option value="attached">Attached</option>
          </select>
        </div>
      </Card>
      {files.isLoading ? <LoadingState /> : null}
      {files.isError ? <ApiErrorState message={files.error.message} /> : null}
      {files.data && files.data.files.length === 0 ? <EmptyState title="No files" description="Upload images here, then reuse their public URLs in Markdown editors." /> : null}
      <div className={compact ? 'grid gap-3 md:grid-cols-2' : 'grid gap-4 md:grid-cols-3 xl:grid-cols-4'}>
        {files.data?.files.map((file) => (
          <Card className="grid gap-3" key={file.id}>
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img alt={file.originalName} className="aspect-video w-full rounded-xl border border-[var(--border)] object-cover" src={file.publicUrl} />
            <div>
              <p className="truncate text-sm font-medium">{file.originalName}</p>
              <p className="text-xs text-[var(--muted-foreground)]">{formatBytes(file.byteSize)} · {file.status}</p>
            </div>
            <div className="flex flex-wrap gap-2">
              {onSelect ? <Button className="px-3 py-1 text-xs" onClick={() => onSelect(file)} type="button">Insert</Button> : null}
              <Button className="bg-red-600 px-3 py-1 text-xs text-white" disabled={deleteFile.isPending} onClick={() => deleteFile.mutate(file.id)} type="button">Delete</Button>
              <a className="rounded-full border border-[var(--border)] px-3 py-1 text-xs" href={file.publicUrl} rel="noreferrer" target="_blank">Open</a>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}

function formatBytes(value: number) {
  if (value < 1024) return `${value} B`;
  if (value < 1024 * 1024) return `${(value / 1024).toFixed(1)} KB`;
  return `${(value / 1024 / 1024).toFixed(1)} MB`;
}
