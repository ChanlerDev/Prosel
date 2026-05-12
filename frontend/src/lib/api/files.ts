import { readAuthState } from '@/lib/auth/store';
import { env } from '@/lib/env';
import { ApiClientError, request } from '@/lib/api/client';
import type { ApiResponse, PaginationMeta } from '@/types/api';
import type { AdminFileListParams, FileAsset, FileListResult } from '@/types/file';

function queryString(params: object) {
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== '') {
      search.set(key, String(value));
    }
  }
  const value = search.toString();
  return value ? `?${value}` : '';
}

function emptyMeta(): PaginationMeta {
  return { page: 1, perPage: 20, total: 0, totalPages: 0 };
}

function accessToken() {
  const auth = readAuthState();
  if (!auth.accessToken) {
    throw new Error('Authentication required');
  }
  return auth.accessToken;
}

async function listWithMeta(path: string): Promise<FileListResult> {
  const response = await request<ApiResponse<FileAsset[]>>(path, { cache: 'no-store', accessToken: accessToken(), rawResponse: true });
  return { files: response.data ?? [], meta: response.meta ?? emptyMeta() };
}

async function upload(file: File) {
  const form = new FormData();
  form.set('file', file);
  const response = await fetch(`${env.apiBaseUrl}/admin/files/upload`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${accessToken()}` },
    body: form,
  });
  const body = (await response.json()) as ApiResponse<FileAsset>;
  if (!response.ok || body.error || !body.data) {
    throw new ApiClientError(body.error?.message ?? 'Upload failed', body.error?.code ?? 'REQUEST_FAILED', response.status);
  }
  return body.data;
}

export const filesApi = {
  admin: {
    list: (params: AdminFileListParams = {}) => listWithMeta(`/admin/files${queryString(params)}`),
    upload,
    attachRef: (id: string, refType: string, refId: string) => request<{ ok: boolean }>(`/admin/files/${id}/ref`, { method: 'PATCH', body: JSON.stringify({ refType, refId }), cache: 'no-store', accessToken: accessToken() }),
    delete: (id: string) => request<{ ok: boolean }>(`/admin/files/${id}`, { method: 'DELETE', cache: 'no-store', accessToken: accessToken() }),
  },
};
