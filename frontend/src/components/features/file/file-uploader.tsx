'use client';

import { ChangeEvent } from 'react';

import { ApiErrorState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { useUploadFile } from '@/lib/files/hooks';
import type { FileAsset } from '@/types/file';

export function FileUploader({ onUploaded }: { onUploaded?: (file: FileAsset) => void }) {
  const upload = useUploadFile();

  return (
    <div className="grid gap-3">
      <label className="grid gap-2 text-sm">
        Upload image
        <input accept="image/jpeg,image/png,image/gif,image/webp" className="text-sm" disabled={upload.isPending} onChange={handleChange} type="file" />
      </label>
      <p className="text-xs text-[var(--muted-foreground)]">Supports JPG, PNG, GIF, and WebP up to the backend upload limit.</p>
      {upload.isError ? <ApiErrorState message={upload.error.message} /> : null}
      {upload.isPending ? <Button disabled type="button">Uploading...</Button> : null}
    </div>
  );

  function handleChange(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    if (!file) return;
    upload.mutate(file, {
      onSuccess: (asset) => {
        onUploaded?.(asset);
        event.target.value = '';
      },
    });
  }
}
