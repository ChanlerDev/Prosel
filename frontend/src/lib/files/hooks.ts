'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import { filesApi } from '@/lib/api/files';
import type { AdminFileListParams } from '@/types/file';

export const fileKeys = {
  adminList: (params: AdminFileListParams) => ['admin', 'files', params] as const,
};

export function useAdminFiles(params: AdminFileListParams) {
  return useQuery({
    queryKey: fileKeys.adminList(params),
    queryFn: () => filesApi.admin.list(params),
  });
}

export function useUploadFile() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (file: File) => filesApi.admin.upload(file),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['admin', 'files'] });
    },
  });
}

export function useAttachFileRef() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, refType, refId }: { id: string; refType: string; refId: string }) => filesApi.admin.attachRef(id, refType, refId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['admin', 'files'] });
    },
  });
}

export function useDeleteFile() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => filesApi.admin.delete(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['admin', 'files'] });
    },
  });
}
