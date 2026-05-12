'use client';

import { useState } from 'react';

import { MediaLibrary } from '@/components/features/file/media-library';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { useAttachFileRef } from '@/lib/files/hooks';
import type { FileAsset } from '@/types/file';

export function MediaPicker({ onSelect, refType, refId }: { onSelect: (file: FileAsset) => void; refType?: string; refId?: string }) {
  const [open, setOpen] = useState(false);
  const attachRef = useAttachFileRef();

  return (
    <div className="grid gap-3">
      <Button className="w-fit px-3 py-1 text-xs" onClick={() => setOpen((current) => !current)} type="button">
        {open ? 'Close media library' : 'Insert image'}
      </Button>
      {open ? (
        <Card className="border-dashed">
          <MediaLibrary compact onSelect={(file) => { onSelect(file); if (refType && refId) attachRef.mutate({ id: file.id, refType, refId }); setOpen(false); }} />
        </Card>
      ) : null}
    </div>
  );
}
